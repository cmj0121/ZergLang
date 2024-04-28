from zergb.lexer import Lexer
from zergb.lexer import TokenType


class TestLexerByToken:
    def setup_method(self):
        self.lexer = Lexer().lexer_by_tokens

    def test_empty(self):
        src = ''
        tokens = list(self.lexer(src))

        assert len(tokens) == 0

    def test_simple_space(self):
        src = ' '
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].raw == ' '
        assert tokens[0].type == TokenType.SPACE

    def test_multiple_spaces(self):
        src = ' \t  \t '
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].type == TokenType.SPACE

    def test_simple_word(self):
        src = 'variable'
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].raw == 'variable'

    def test_simple_comment(self):
        src = '// This is a comment'
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].raw == '// This is a comment'
        assert tokens[0].type == TokenType.COMMENT

    def test_simple_str(self):
        src = '"Hello World"'
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].raw == '"Hello World"'
        assert tokens[0].type == TokenType.STRING

    def test_newline(self):
        src = '1\n2'
        tokens = list(self.lexer(src))

        assert len(tokens) == 3
        assert tokens[0].raw == '1'
        assert tokens[1].raw == '\n'
        assert tokens[1].type == TokenType.NEWLINE
        assert tokens[2].raw == '2'

    def test_simple_expression(self):
        src = '1 + 2'
        tokens = list(self.lexer(src))

        assert len(tokens) == 5
        assert tokens[0].raw == '1'
        assert tokens[1].type == TokenType.SPACE
        assert tokens[2].raw == '+'
        assert tokens[3].type == TokenType.SPACE
        assert tokens[4].raw == '2'

    def test_expression(self):
        src = '1+2'
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].raw == '1+2'


class TestLexerByOperator:
    def setup_method(self):
        self.lexer = Lexer().lexer

    def test_simple_expr(self):
        src = "1 + 2"
        tokens = list(self.lexer(src))

        assert len(tokens) == 3
        assert tokens[0].raw == '1'
        assert tokens[1].raw == '+'
        assert tokens[1].type == TokenType.ADD
        assert tokens[2].raw == '2'

    def test_expression(self):
        src = "1+2"
        tokens = list(self.lexer(src))

        assert len(tokens) == 3
        assert tokens[0].raw == '1'
        assert tokens[1].raw == '+'
        assert tokens[1].type == TokenType.ADD
        assert tokens[2].raw == '2'

    def test_arithmetic(self):
        src = "-1 + ~2 * +3 / 4 % 5"
        tokens = list(self.lexer(src))

        assert len(tokens) == 12
        assert tokens[0].raw == '-'
        assert tokens[0].type == TokenType.SUB
        assert tokens[1].raw == '1'
        assert tokens[2].raw == '+'
        assert tokens[2].type == TokenType.ADD
        assert tokens[3].raw == '~'
        assert tokens[3].type == TokenType.NEG
        assert tokens[4].raw == '2'
        assert tokens[5].raw == '*'
        assert tokens[5].type == TokenType.MUL
        assert tokens[6].raw == '+'
        assert tokens[6].type == TokenType.ADD
        assert tokens[7].raw == '3'
        assert tokens[8].raw == '/'
        assert tokens[8].type == TokenType.DIV
        assert tokens[9].raw == '4'
        assert tokens[10].raw == '%'
        assert tokens[10].type == TokenType.MOD
        assert tokens[11].raw == '5'

    def test_combine_op(self):
        src = "++"
        tokens = list(self.lexer(src))

        assert len(tokens) == 1
        assert tokens[0].raw == '++'
        assert tokens[0].type == TokenType.INC

    def test_func(self):
        src = "main()"
        tokens = list(self.lexer(src))

        assert len(tokens) == 3
        assert tokens[0].raw == 'main'
        assert tokens[1].raw == '('
        assert tokens[2].raw == ')'
