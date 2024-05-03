from zergb.lexer import TokenType
from zergb.parser import Parser


class TestParser:
    def test_parse_nop(self):
        src = 'nop'
        parser = Parser()
        node = parser.parse(src)

        assert len(node.childs) == 1
        assert node.childs[0].token.type == TokenType.NOP

    def test_parse_fn(self):
        src = 'fn main() { }'
        parser = Parser()
        node = parser.parse(src)

        assert len(node.childs) == 1

        fn = node.childs[0]
        assert fn.token.type == TokenType.FN
        assert len(fn.childs) == 2

        name = fn.childs[0]
        assert name.token.type == TokenType.NAME
        assert name.token.raw == 'main'
