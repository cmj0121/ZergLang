from __future__ import annotations

from zergb.lexer import Lexer
from zergb.lexer import Token
from zergb.lexer import TokenType


class ASTNode:
    def __init__(self, token: Token | None = None):
        self._token = token or Token('.', tt=TokenType.ROOT)
        self._parent = None
        self._childs = []

    def __add__(self, other: ASTNode) -> ASTNode:
        assert isinstance(other, ASTNode)

        self._childs.append(other)
        other._parent = self

        return self

    def __contains__(self, item: ASTNode) -> bool:
        return item in self._childs

    def __str__(self, indent=0) -> str:
        trees = []
        if self.is_root:
            trees.append(self.token.__str__())
        else:
            ch = '├─' if not self.is_latest else '└─'
            trees.append(f'{" " * indent}{ch}  {self.token}')

        for child in self.childs:
            trees.append(child.__str__(indent + 4))

        return '\n'.join(trees)

    @property
    def token(self) -> Token:
        return self._token

    @property
    def parent(self) -> ASTNode | None:
        return self._parent

    @property
    def childs(self) -> list[ASTNode]:
        return self._childs

    @property
    def is_root(self) -> bool:
        return self._parent is None

    @property
    def is_latest(self) -> bool:
        return self.is_root or self.parent.childs[-1] == self


class Parser(Lexer):
    """Parse and build the AST from the tokens"""
    def parse(self, src: str) -> ASTNode:
        self._tokens = self.lexer(src)
        self._token_stack = []
        return self._parse_source()

    @property
    def next_token(self) -> Token:
        if self._token_stack:
            return self._token_stack.pop()
        return next(self._tokens)

    @next_token.setter
    def next_token(self, v: Token):
        self._token_stack.append(v)

    def _parse_source(self) -> ASTNode:
        root = ASTNode()
        prev = self.next_token
        match prev.type:
            case TokenType.RBRACE:
                self.next_token = prev
            case _:
                root + self._parse_block(prev)
        return root

    def _parse_scope(self) -> ASTNode:
        prev = self.next_token
        if prev.type != TokenType.LBRACE:
            raise ValueError(f'expecting {{, but got {prev=}')

        node = self._parse_source()
        prev = self.next_token
        if prev.type != TokenType.RBRACE:
            raise ValueError(f'expecting }} but got {prev=}')

        return node

    def _parse_block(self, prev: Token) -> ASTNode:
        match prev.type:
            case TokenType.NOP:
                node = ASTNode(prev)
            case TokenType.FN:
                node = self._parse_func_stmt(prev)
            case _:
                raise ValueError(f'Unexpected token: {prev}')
        return node

    def _parse_func_stmt(self, prev: Token) -> ASTNode:
        assert prev.type == TokenType.FN
        node = ASTNode(prev)

        node + self._parse_func_head()
        node + self._parse_scope()
        return node

    def _parse_func_head(self) -> ASTNode:
        name = self.next_token
        if name.type != TokenType.NAME:
            raise ValueError(f'function should follow with name: {name=}')
        node = ASTNode(name)

        # handle the function arguments
        prev = self.next_token
        if prev.type != TokenType.LPARENTHESES:
            raise ValueError(f'function should follow with (): {prev=}')
        prev = self.next_token
        match prev.type:
            case TokenType.RPARENTHESES:
                pass
            case _:
                node + self._parse_func_args(prev)

        # handle the function type hint
        prev = self.next_token
        match prev.type:
            case TokenType.ARROW:
                node + self._parse_type_hint(prev)
            case _:
                self.next_token = prev

        return node

    def _parse_func_args(self, prev: Token) -> ASTNode:
        raise NotImplementedError
