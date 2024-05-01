#! /usr/bin/env python
from llvmlite import binding as llvm
from llvmlite import ir
from llvmlite.ir import Module
from zergb.lexer import TokenType
from zergb.parser import ASTNode
from zergb.parser import Parser


class ZergBootstrap(Parser):
    def compile(self, source: str) -> bytes:
        '''compile source code from filepath to object code'''
        ir = self._source_to_ir(source)
        return self._ir_to_obj(ir)

    def _source_to_ir(self, source: str) -> str:
        '''compile source code from filepath to LLVM IR'''
        llvm.initialize()
        llvm.initialize_native_target()
        llvm.initialize_native_asmprinter()

        with open(source) as f:
            ast = self.parse(f.read())
            module = Module(name=source)
            self._compile_ast(ast, module)

            # Create a LLVM module object from the IR
            mod = llvm.parse_assembly(str(module))
            mod.verify()

            return mod

    def _ir_to_obj(self, ir: str) -> bytes:
        '''compile LLVM IR to object code'''
        target = llvm.Target.from_default_triple()
        target_machine = target.create_target_machine()

        return target_machine.emit_object(ir)

    def _compile_ast(self, ast: ASTNode, module: Module | None = None, builder: ir.IRBuilder | None = None):
        match ast.token.type:
            case TokenType.ROOT:
                for child in ast.childs:
                    self._compile_ast(child, module, builder)
            case TokenType.FN:
                name, body = ast.childs
                assert name.token.type == TokenType.NAME
                assert body.token.type == TokenType.ROOT

                fnty = ir.FunctionType(ir.VoidType(), [])
                func = ir.Function(module, fnty, name=name.token.raw)

                block = func.append_basic_block(name='entry')
                builder = ir.IRBuilder(block)

                self._compile_ast(body, module, builder)
                builder.ret_void()
            case TokenType.NOP:
                assert builder is not None, 'NOP outside of function'
            case _:
                raise NotImplementedError(f'Unknown token type: {ast.token.type}')
