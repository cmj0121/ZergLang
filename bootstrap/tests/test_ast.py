from zergb.lexer import Token
from zergb.lexer import TokenType
from zergb.parser import ASTNode


class TestASTNode:
    def test_dummy_node(self):
        token = Token(' ', tt=TokenType.SPACE)
        node = ASTNode(token)

        assert node.token == token
        assert node.parent is None
        assert node.childs == []
        assert str(node) == '.\n└─  [SPACE]'

    def test_add_child(self):
        token = Token(' ', tt=TokenType.SPACE)
        node = ASTNode(token)

        child = ASTNode(Token('variable'))
        node + child

        assert node.childs == [child]
        assert child.parent == node
        assert child in node
        assert str(node) == '.\n└─  [SPACE]\n    └─  variable'

    def test_arithmetic(self):
        token = Token('+')
        node = ASTNode(token)

        node + ASTNode(Token('1'))
        node + ASTNode(Token('2'))

        assert str(node) == '.\n└─  +\n    ├─  1\n    └─  2'
