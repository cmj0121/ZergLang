from __future__ import annotations

from zergb.lexer import Token


class ASTNode:
    def __init__(self, token: Token):
        self._token = token
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
        if indent == 0:
            trees.append('.')

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
