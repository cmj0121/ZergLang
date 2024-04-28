import enum
from typing import Generator


class TokenType(enum.Enum):
    """The classified token of the zergb"""
    UNKNOWN = enum.auto()
    NEWLINE = enum.auto()
    COMMENT = enum.auto()

    INDENT = enum.auto()
    DEDENT = enum.auto()
    SPACE = enum.auto()

    STRING = enum.auto()


class Token:
    def __init__(self, raw: str, tt: TokenType = TokenType.UNKNOWN):
        self._raw = raw
        self._type = tt

    def __repr__(self):
        return self._raw

    @property
    def raw(self) -> str:
        return self._raw

    @property
    def type(self) -> TokenType:
        return self._type


class Lexer:
    """The lexer class for the zergb"""
    def lexer(self, src: str) -> Generator[Token, None, None]:
        '''tokenize the source code and return the generator of tokens'''
        yield from self._lexer(src)

    def _lexer(self, src: str) -> Generator[Token, None, None]:
        '''tokenize the source code by several lexers'''
        yield from self.lexer_by_tokens(src)

    def lexer_by_tokens(self, src) -> Generator[Token, None, None]:
        '''
        tokenize the source code with tokens.

        In this stage return the tokens and only NEWLINE, COMMENT, STRING and
        SPACE are classified.
        '''
        index = 0

        while index < len(src):
            tt = src[index]

            match tt:
                case '\n':
                    yield Token(tt, TokenType.NEWLINE)
                case '/':
                    if index + 1 < len(src) and src[index + 1] == '/':
                        stop = index + 2
                        while stop < len(src) and src[stop] != '\n':
                            stop += 1
                        yield Token(src[index:stop], TokenType.COMMENT)
                        index = stop
                    else:
                        yield Token(tt)
                case ' ' | '\t':
                    stop = index + 1
                    while stop < len(src) and src[stop] in ' \t':
                        stop += 1
                    yield Token(src[index:stop], TokenType.SPACE)
                    index = stop - 1
                case '"':
                    stop = index + 1
                    while stop < len(src) and src[stop] != '"':
                        stop += 1
                    yield Token(src[index:stop + 1], TokenType.STRING)
                    index = stop
                case _:
                    stop = index
                    while stop < len(src) and src[stop] not in ' \t\n':
                        stop += 1
                    yield Token(src[index:stop])
                    index = stop - 1

            index += 1
