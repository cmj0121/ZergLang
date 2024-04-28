import enum
from typing import Generator


@enum.unique
class TokenType(enum.Enum):
    """The classified token of the zergb"""
    UNKNOWN = enum.auto()
    NEWLINE = enum.auto()
    COMMENT = enum.auto()

    INDENT = enum.auto()
    DEDENT = enum.auto()
    SPACE = enum.auto()

    STRING = enum.auto()

    # known operators
    ADD = '+'
    SUB = '-'
    MUL = '*'
    DIV = '/'
    MOD = '%'
    NEG = '~'
    INC = '++'
    DEC = '--'
    LT = '<'
    GT = '>'
    AND = '&'
    OR = '|'
    NOT = '!'
    XOR = '^'
    LSHIFT = '<<'
    RSHIFT = '>>'
    LPARENTHESES = '('
    RPARENTHESES = ')'
    LBRACE = '{'
    RBRACE = '}'
    LBRACKET = '['
    RBRACKET = ']'

    @staticmethod
    def yield_tokens(raw: str) -> Generator['TokenType', None, None]:
        '''yield the token type from the raw string'''
        try:
            yield Token(raw, tt=TokenType(raw))
        except ValueError:
            first, remains = raw[0], raw[1:]
            yield Token(first, tt=TokenType(first))
            yield from TokenType.yield_tokens(remains)


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
    OPERATORS = '+-*/%<>&|!^~(){}[]'

    """The lexer class for the zergb"""
    def lexer(self, src: str) -> Generator[Token, None, None]:
        '''tokenize the source code and return the generator of tokens'''
        yield from self._lexer(src)

    def _lexer(self, src: str) -> Generator[Token, None, None]:
        '''tokenize the source code by several lexers'''
        base = self.lexer_by_tokens(src)
        base = self.lexer_extract_operator(base)
        # the latest lexer and remove the unuseful tokens
        base = self.lexer_remove_unuseful(base)

        yield from base

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

    def lexer_extract_operator(self, tokens) -> Generator[Token, None, None]:
        """extract the operator from the unknown tokens"""
        for token in tokens:
            match token.type:
                case TokenType.UNKNOWN:
                    # only process the token with unknown type
                    remains, operators = '', ''
                    for tt in token.raw:
                        if tt in self.OPERATORS:
                            if remains:
                                yield Token(remains)
                                remains = ''
                            operators += tt
                        else:
                            if operators:
                                yield from TokenType.yield_tokens(operators)
                                operators = ''
                            remains += tt

                    if remains:
                        yield Token(remains)
                    if operators:
                        yield from TokenType.yield_tokens(operators)
                case _:
                    yield token

    def lexer_remove_unuseful(self, tokens) -> Generator[Token, None, None]:
        for token in tokens:
            match token.type:
                case TokenType.SPACE | TokenType.COMMENT | TokenType.NEWLINE:
                    pass
                case _:
                    yield token
