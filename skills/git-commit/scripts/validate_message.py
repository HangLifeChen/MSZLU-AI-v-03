#!/usr/bin/env python3
"""
Conventional Commits 提交信息验证脚本

功能：
1. 验证提交信息是否符合 Conventional Commits 规范
2. 检查类型、范围、描述、正文、页脚的格式
3. 提供详细的错误提示和修复建议

使用方式：
    python validate_message.py "feat(auth): add login support"
    python validate_message.py --file .git/COMMIT_EDITMSG
    echo "feat: add feature" | python validate_message.py --stdin
"""

import sys
import re
import argparse
from typing import List, Tuple, Optional
from enum import Enum

# 允许的类型（可根据项目需要调整）
VALID_TYPES = {
    'feat', 'fix', 'docs', 'style', 'refactor', 
    'perf', 'test', 'chore', 'ci', 'build', 'revert'
}

# 最大长度限制
MAX_HEADER_LENGTH = 72  # GitHub建议
MAX_SUBJECT_LENGTH = 50  # 主题建议更短
MAX_BODY_LINE_LENGTH = 72


class ValidationError(Enum):
    """验证错误类型"""
    EMPTY_MESSAGE = "提交信息不能为空"
    INVALID_FORMAT = "格式错误：必须符合 <type>[scope]: <subject> 格式"
    INVALID_TYPE = "无效的类型 '{}'，允许的类型: {}"
    EMPTY_SCOPE = "范围(scope)不能为空（如果使用括号）"
    UPPERCASE_SUBJECT = "描述(subject)必须以**小写字母**开头"
    SUBJECT_ENDS_WITH_PERIOD = "描述(subject)**不能以句号结尾**"
    SUBJECT_TOO_LONG = "描述过长 ({} > {} 字符)，建议控制在{}字符以内"
    NON_IMPERATIVE = "描述建议使用祈使句（如 'add' 而非 'added' 或 'adds'）"
    LINE_TOO_LONG = "第{}行过长 ({} > {} 字符)"
    INVALID_BREAKING_CHANGE = "破坏性变更格式错误，应以 BREAKING CHANGE: 开头"
    MISSING_BLANK_LINE = "标题与正文之间需要空一行"
    BODY_EMPTY_LINE = "正文段落之间建议空行分隔，便于阅读"


class CommitValidator:
    """提交信息验证器"""
    
    def __init__(self, message: str):
        self.message = message.strip()
        self.errors: List[Tuple[ValidationError, str]] = []
        self.warnings: List[Tuple[ValidationError, str]] = []
        
    def validate(self) -> bool:
        """执行完整验证"""
        if not self.message:
            self.errors.append((ValidationError.EMPTY_MESSAGE, ""))
            return False
            
        lines = self.message.split('\n')
        
        # 1. 验证标题行（第一行）
        if not self._validate_header(lines[0]):
            return False
            
        # 2. 验证正文结构（如果有）
        if len(lines) > 1:
            self._validate_body_structure(lines)
            
        # 3. 验证页脚（破坏性变更、Issue引用等）
        self._validate_footer(lines)
        
        return len(self.errors) == 0
    
    def _validate_header(self, header: str) -> bool:
        """验证标题行格式：<type>[optional scope]: <description>"""
        # 匹配模式：type(scope): subject 或 type!: subject 或 type(scope)!: subject
        pattern = r'^(?P<type>\w+)(?:\((?P<scope>[^)]+)\))?(?P<breaking>!)?: (?P<subject>.+)$'
        match = re.match(pattern, header)
        
        if not match:
            self.errors.append((ValidationError.INVALID_FORMAT, header))
            return False
            
        groups = match.groupdict()
        commit_type = groups['type']
        scope = groups['scope']
        is_breaking = groups['breaking'] == '!'
        subject = groups['subject']
        
        # 验证类型
        if commit_type not in VALID_TYPES:
            self.errors.append((
                ValidationError.INVALID_TYPE, 
                f"{commit_type}, {', '.join(sorted(VALID_TYPES))}"
            ))
            return False
        
        # 验证范围（如果存在）
        if scope is not None and not scope.strip():
            self.errors.append((ValidationError.EMPTY_SCOPE, ""))
            return False
            
        # 验证描述
        self._validate_subject(subject, is_breaking)
        
        return len(self.errors) == 0
    
    def _validate_subject(self, subject: str, is_breaking: bool) -> None:
        """验证描述部分"""
        # 检查首字母大写
        if subject and subject[0].isupper():
            self.errors.append((ValidationError.UPPERCASE_SUBJECT, subject))
        
        # 检查以句号结尾（允许使用...表示省略，但不允许.结尾）
        if subject.rstrip().endswith('.') and not subject.rstrip().endswith('...'):
            self.errors.append((ValidationError.SUBJECT_ENDS_WITH_PERIOD, subject))
        
        # 检查长度（如果是破坏性变更，允许更长一些，因为有!标记）
        max_len = MAX_SUBJECT_LENGTH if not is_breaking else MAX_SUBJECT_LENGTH - 1
        actual_len = len(subject)
        if actual_len > MAX_HEADER_LENGTH:
            self.errors.append((
                ValidationError.SUBJECT_TOO_LONG, 
                f"{actual_len}, {MAX_HEADER_LENGTH}, {max_len}"
            ))
        elif actual_len > max_len:
            self.warnings.append((
                ValidationError.SUBJECT_TOO_LONG, 
                f"{actual_len}, {MAX_HEADER_LENGTH}, {max_len}"
            ))
        
        # 检查是否使用过去时（简单启发式检查）
        past_tense_words = ['added', 'fixed', 'updated', 'removed', 'changed', 'implemented', 'created', 'modified']
        first_word = subject.split()[0].lower().rstrip(':,;.')
        if first_word in past_tense_words:
            self.warnings.append((ValidationError.NON_IMPERATIVE, first_word))
    
    def _validate_body_structure(self, lines: List[str]) -> None:
        """验证正文结构"""
        # 检查标题和正文之间是否有空行
        if len(lines) > 1 and lines[1].strip():
            self.errors.append((ValidationError.MISSING_BLANK_LINE, ""))
        
        # 检查每行长度
        for i, line in enumerate(lines[2:], start=3):  # 从第3行开始（0-indexed）
            if len(line) > MAX_BODY_LINE_LENGTH:
                self.errors.append((
                    ValidationError.LINE_TOO_LONG, 
                    f"{i}, {len(line)}, {MAX_BODY_LINE_LENGTH}"
                ))
    
    def _validate_footer(self, lines: List[str]) -> None:
        """验证页脚格式（破坏性变更、Issue引用等）"""
        in_footer = False
        for i, line in enumerate(lines[1:], start=2):
            # 检测页脚开始（通常是BREAKING CHANGE或Issue引用）
            if re.match(r'^(BREAKING CHANGE|Closes|Fixes|Refs|Relates to)[\s:#]', line):
                in_footer = True
            
            if in_footer:
                # 验证BREAKING CHANGE格式
                if line.startswith('BREAKING CHANGE') and not re.match(r'^BREAKING CHANGE:', line):
                    self.errors.append((ValidationError.INVALID_BREAKING_CHANGE, line))
    
    def get_report(self) -> str:
        """生成验证报告"""
        lines = []
        
        if self.errors:
            lines.append("❌ 验证失败，发现以下错误：")
            for error, detail in self.errors:
                msg = error.value
                if detail:
                    msg = msg.format(*detail.split(', ')) if '{}' in msg else f"{msg} ({detail})"
                lines.append(f"   - {msg}")
        
        if self.warnings:
            lines.append("\n⚠️  警告（建议修复）：")
            for warning, detail in self.warnings:
                msg = warning.value
                if detail:
                    msg = msg.format(*detail.split(', ')) if '{}' in msg else f"{msg} ({detail})"
                lines.append(f"   - {msg}")
        
        if not self.errors and not self.warnings:
            lines.append("✅ 提交信息符合 Conventional Commits 规范")
            
        return '\n'.join(lines)
    
    def is_valid(self) -> bool:
        """返回是否验证通过（无错误）"""
        return len(self.errors) == 0


def read_from_file(filepath: str) -> str:
    """从文件读取提交信息，跳过注释行"""
    with open(filepath, 'r', encoding='utf-8') as f:
        lines = []
        for line in f:
            # 跳过Git注释行和空行
            line = line.rstrip()
            if not line.startswith('#') and line.strip():
                lines.append(line)
        return '\n'.join(lines)


def read_from_stdin() -> str:
    """从标准输入读取"""
    return sys.stdin.read()


def main():
    parser = argparse.ArgumentParser(
        description='验证Git提交信息是否符合 Conventional Commits 规范',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
示例:
  %(prog)s "feat(auth): add OAuth login"
  %(prog)s --file .git/COMMIT_EDITMSG
  git log -1 --pretty=%%B | %(prog)s --stdin
        '''
    )
    parser.add_argument('message', nargs='?', help='直接传入提交信息')
    parser.add_argument('-f', '--file', help='从文件读取提交信息')
    parser.add_argument('-s', '--stdin', action='store_true', help='从标准输入读取')
    parser.add_argument('-q', '--quiet', action='store_true', help='静默模式，只返回退出码')
    parser.add_argument('--strict', action='store_true', help='严格模式（警告视为错误）')
    
    args = parser.parse_args()
    
    # 获取提交信息
    message = ""
    if args.stdin:
        message = read_from_stdin()
    elif args.file:
        try:
            message = read_from_file(args.file)
        except FileNotFoundError:
            print(f"错误：找不到文件 {args.file}")
            sys.exit(2)
    elif args.message:
        message = args.message
    else:
        # 尝试从 .git/COMMIT_EDITMSG 读取（Git hook场景）
        if os.path.exists('.git/COMMIT_EDITMSG'):
            message = read_from_file('.git/COMMIT_EDITMSG')
        else:
            parser.print_help()
            sys.exit(2)
    
    # 验证
    validator = CommitValidator(message)
    is_valid = validator.validate()
    
    # 严格模式下，警告也视为错误
    if args.strict and validator.warnings:
        is_valid = False
    
    # 输出结果
    if not args.quiet:
        print(validator.get_report())
        print(f"\n提交信息预览:\n{'-'*40}")
        print(message)
        print('-'*40)
    
    sys.exit(0 if is_valid else 1)


if __name__ == "__main__":
    main()