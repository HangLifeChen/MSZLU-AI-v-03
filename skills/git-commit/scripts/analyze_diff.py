#!/usr/bin/env python3
"""
Git Diff 分析脚本
分析工作区变更，推断提交类型、影响范围和变更摘要
"""

import subprocess
import json
import sys
import re
import os
from typing import List, Dict, Optional, Tuple


class GitDiffAnalyzer:
    def __init__(self):
        self.repo_root = self._get_repo_root()
        
    def _get_repo_root(self) -> str:
        """获取Git仓库根目录"""
        try:
            result = subprocess.run(
                ['git', 'rev-parse', '--show-toplevel'],
                capture_output=True, text=True, check=True
            )
            return result.stdout.strip()
        except subprocess.CalledProcessError:
            print(json.dumps({"error": "不在Git仓库中或Git命令执行失败"}))
            sys.exit(1)
        except FileNotFoundError:
            print(json.dumps({"error": "未找到Git命令，请确保Git已安装"}))
            sys.exit(1)
    
    def _run_git(self, args: List[str]) -> Tuple[str, str, int]:
        """运行Git命令并返回结果"""
        try:
            result = subprocess.run(
                ['git'] + args,
                capture_output=True, text=True,
                cwd=self.repo_root
            )
            return result.stdout, result.stderr, result.returncode
        except Exception as e:
            return "", str(e), 1
    
    def get_changed_files(self, cached: bool = False) -> List[Dict]:
        """
        获取变更文件列表
        
        Args:
            cached: 是否只查看已暂存的变更
            
        Returns:
            文件变更列表，包含状态、路径、重命名信息等
        """
        cmd = ['diff', '--name-status']
        if cached:
            cmd.append('--cached')
            
        stdout, stderr, code = self._run_git(cmd)
        
        if code != 0:
            return [{"error": f"获取变更失败: {stderr}"}]
        
        files = []
        lines = [l.strip() for l in stdout.split('\n') if l.strip()]
        
        i = 0
        while i < len(lines):
            line = lines[i]
            parts = line.split('\t')
            status = parts[0][0]  # 获取首字母 (A/M/D/R/C)
            
            file_info = {
                'status': status,
                'status_desc': self._get_status_desc(status),
                'path': parts[-1],
                'is_staged': cached
            }
            
            # 处理重命名(R)和复制(C) - 格式: R100\told_path\tnew_path
            if status in ['R', 'C'] and len(parts) >= 3:
                file_info['old_path'] = parts[1]
                file_info['path'] = parts[2]
                file_info['similarity'] = parts[0][1:]  # 相似度百分比
            
            # 获取该文件的diff统计（添加/删除行数）
            stat = self._get_file_diff_stat(parts[-1], cached)
            file_info.update(stat)
            
            files.append(file_info)
            i += 1
            
        return files
    
    def _get_status_desc(self, status: str) -> str:
        """获取状态描述"""
        status_map = {
            'A': 'added',      # 新增
            'M': 'modified',   # 修改
            'D': 'deleted',    # 删除
            'R': 'renamed',    # 重命名
            'C': 'copied',     # 复制
            'U': 'updated',    # 更新但未合并
            '?': 'untracked'   # 未跟踪（使用--name-status不会显示，但备用）
        }
        return status_map.get(status, 'unknown')
    
    def _get_file_diff_stat(self, file_path: str, cached: bool) -> Dict:
        """获取单个文件的diff统计信息"""
        cmd = ['diff', '--numstat']
        if cached:
            cmd.append('--cached')
        cmd.append('--')
        cmd.append(file_path)
        
        stdout, _, _ = self._run_git(cmd)
        if stdout.strip():
            # 格式: <添加行数>\t<删除行数>\t<文件路径>
            parts = stdout.strip().split('\t')
            if len(parts) >= 2:
                try:
                    return {
                        'additions': int(parts[0]) if parts[0] != '-' else 0,
                        'deletions': int(parts[1]) if parts[1] != '-' else 0
                    }
                except ValueError:
                    pass
        return {'additions': 0, 'deletions': 0}
    
    def get_diff_content(self, file_path: str, cached: bool = False) -> str:
        """获取指定文件的详细diff内容"""
        cmd = ['diff']
        if cached:
            cmd.append('--cached')
        cmd.extend(['--', file_path])
        
        stdout, _, _ = self._run_git(cmd)
        return stdout
    
    def analyze_scope(self, file_path: str) -> Optional[str]:
        """
        根据文件路径推断scope（模块/组件）
        
        策略：
        1. 优先返回第一级目录（如果是常见源码目录则返回第二级）
        2. 测试文件返回测试目标模块名
        """
        parts = file_path.split('/')
        if not parts:
            return None
        
        # 常见源码目录映射
        skip_dirs = {'src', 'lib', 'source', 'code', 'pkg'}
        
        # 如果是测试文件，尝试提取被测模块
        if 'test' in parts or 'tests' in parts or file_path.endswith(('_test.go', '.test.ts', '.spec.ts', '.test.js', '.spec.js')):
            # 尝试找到对应的源码目录
            for i, part in enumerate(parts):
                if part in ('test', 'tests'):
                    if i + 1 < len(parts):
                        return parts[i + 1]
                    break
            # 从文件名推断（如 auth_test.go -> auth）
            basename = os.path.basename(file_path)
            for ext in ['_test.go', '.test.ts', '.spec.ts', '.test.js', '.spec.js']:
                if basename.endswith(ext):
                    return basename[:-len(ext)]
        
        # 常规路径分析
        if len(parts) > 1:
            first_dir = parts[0]
            if first_dir in skip_dirs and len(parts) > 2:
                return parts[1]
            elif first_dir.startswith('.'):  # 隐藏目录如.github
                return parts[1] if len(parts) > 1 else first_dir
            else:
                return first_dir
        
        return None
    
    def suggest_type(self, file_info: Dict, diff_content: str = "") -> str:
        """
        根据文件变更内容推断提交类型
        
        Args:
            file_info: 文件信息字典
            diff_content: 文件diff内容（可选，用于更精确分析）
        """
        path = file_info['path'].lower()
        status = file_info['status']
        
        # 基于文件路径的快速判断
        if status == 'D':
            return 'refactor'  # 删除通常是重构
        
        # 文档文件
        if any(path.endswith(ext) for ext in ['.md', '.rst', '.txt', '.doc']) or \
           'docs' in path or 'doc' in path:
            return 'docs'
        
        # 配置文件
        if any(path.endswith(ext) for ext in ['.yml', '.yaml', '.json', '.toml', '.ini', '.cfg']) or \
           path.startswith('.github/') or path.startswith('.ci/'):
            return 'chore'
        
        # 测试文件
        if any(pattern in path for pattern in ['test', 'spec', '__tests__']) or \
           path.endswith(('_test.go', '.test.ts', '.spec.ts')):
            return 'test'
        
        # 基于diff内容分析（如果提供了内容）
        if diff_content:
            diff_lower = diff_content.lower()
            
            # 关键词匹配
            if any(kw in diff_lower for kw in ['fix', 'bug', 'hotfix', 'patch', '修复']):
                return 'fix'
            elif any(kw in diff_lower for kw in ['refactor', 'clean', 'restructure', '重构']):
                return 'refactor'
            elif any(kw in diff_lower for kw in ['perf', 'optimize', 'performance', '加速', '优化']):
                return 'perf'
            elif any(kw in diff_lower for kw in ['add', 'new', 'feature', 'implement', '新增', '添加']):
                return 'feat'
            elif any(kw in diff_lower for kw in ['deprecated', 'obsolete', '弃用']):
                return 'refactor'
        
        # 基于变更统计的启发式判断
        additions = file_info.get('additions', 0)
        deletions = file_info.get('deletions', 0)
        
        # 大量删除通常是重构
        if deletions > additions * 2 and deletions > 50:
            return 'refactor'
        
        # 默认假设为新功能（保守选择，避免误标为fix）
        return 'feat'
    
    def analyze_all(self, cached: bool = False, include_diff: bool = False) -> Dict:
        """
        执行完整分析
        
        Args:
            cached: 是否只分析已暂存变更
            include_diff: 是否包含详细diff内容（注意：可能很大）
        """
        files = self.get_changed_files(cached)
        
        if not files:
            return {
                "summary": {
                    "total_files": 0,
                    "message": "没有检测到变更" + ("（已暂存）" if cached else "（未暂存）")
                },
                "files": []
            }
        
        # 分析每个文件
        analyzed_files = []
        scopes = set()
        types = {}
        
        total_additions = 0
        total_deletions = 0
        
        for f in files:
            if 'error' in f:
                continue
                
            # 推断scope
            scope = self.analyze_scope(f['path'])
            if scope:
                scopes.add(scope)
                f['inferred_scope'] = scope
            
            # 获取diff内容并推断类型（如果内容不大）
            diff_content = ""
            if include_diff or f.get('additions', 0) < 100:  # 小文件才获取内容
                diff_content = self.get_diff_content(f['path'], cached)
            
            commit_type = self.suggest_type(f, diff_content)
            f['suggested_type'] = commit_type
            types[commit_type] = types.get(commit_type, 0) + 1
            
            if include_diff:
                f['diff_sample'] = diff_content[:500] + '...' if len(diff_content) > 500 else diff_content
            
            total_additions += f.get('additions', 0)
            total_deletions += f.get('deletions', 0)
            analyzed_files.append(f)
        
        # 确定主导类型（主要提交类型）
        dominant_type = max(types, key=types.get) if types else 'feat'
        
        # 确定主导scope
        dominant_scope = None
        if scopes:
            # 选择变更文件数最多的scope
            scope_counts = {}
            for f in analyzed_files:
                if 'inferred_scope' in f:
                    s = f['inferred_scope']
                    scope_counts[s] = scope_counts.get(s, 0) + 1
            dominant_scope = max(scope_counts, key=scope_counts.get)
        
        return {
            "summary": {
                "total_files": len(analyzed_files),
                "total_additions": total_additions,
                "total_deletions": total_deletions,
                "suggested_type": dominant_type,
                "suggested_scope": dominant_scope,
                "type_breakdown": types,
                "detected_scopes": sorted(list(scopes)),
                "is_staged": cached
            },
            "files": analyzed_files
        }


def main():
    import argparse
    parser = argparse.ArgumentParser(description='分析Git变更并推断提交信息')
    parser.add_argument('--cached', action='store_true', 
                       help='仅分析已暂存的变更(git add后的)')
    parser.add_argument('--include-diff', action='store_true',
                       help='包含详细diff内容（可能输出很大）')
    parser.add_argument('--format', choices=['json', 'pretty'], default='json',
                       help='输出格式')
    
    args = parser.parse_args()
    
    analyzer = GitDiffAnalyzer()
    result = analyzer.analyze_all(cached=args.cached, include_diff=args.include_diff)
    
    if args.format == 'json':
        print(json.dumps(result, indent=2, ensure_ascii=False))
    else:
        # 友好格式输出
        s = result['summary']
        print(f"📊 变更统计:")
        print(f"   文件数: {s['total_files']}")
        print(f"   添加: +{s['total_additions']}, 删除: -{s['total_deletions']}")
        print(f"\n📝 建议提交类型: {s['suggested_type']}")
        if s['suggested_scope']:
            print(f"📦 建议Scope: {s['suggested_scope']}")
        print(f"\n📁 文件明细:")
        for f in result['files']:
            status_emoji = {'A': '➕', 'M': '✏️', 'D': '🗑️', 'R': '📝'}.get(f['status'], '•')
            scope_info = f"[{f['inferred_scope']}]" if 'inferred_scope' in f else ''
            type_info = f"({f['suggested_type']})" if 'suggested_type' in f else ''
            print(f"   {status_emoji} {f['path']} {scope_info} {type_info} +{f.get('additions',0)}/-{f.get('deletions',0)}")


if __name__ == "__main__":
    main() 