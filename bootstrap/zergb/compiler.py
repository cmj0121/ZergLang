#! /usr/bin/env python
import llvmlite.binding as llvm
import llvmlite.ir


class ZergBootstrap:
    def compile(self, src: str) -> bytes:
        '''compile source code to object code'''
        ir = self._source_to_ir(src)
        return self._ir_to_obj(ir)

    def _source_to_ir(self, src: str) -> str:
        '''compile source code to LLVM IR'''
        llvm.initialize()
        llvm.initialize_native_target()
        llvm.initialize_native_asmprinter()

        return self._compile(src)

    def _ir_to_obj(self, ir: str) -> bytes:
        '''compile LLVM IR to object code'''
        target = llvm.Target.from_default_triple()
        target_machine = target.create_target_machine()

        return target_machine.emit_object(ir)

    def _compile(self, src: str) -> str:
        module = llvmlite.ir.Module(name=__file__)
        fnty = llvmlite.ir.FunctionType(llvmlite.ir.IntType(32), [])
        func = llvmlite.ir.Function(module, fnty, name='main')

        builder = llvmlite.ir.IRBuilder(func.append_basic_block('entry'))
        builder.ret(llvmlite.ir.Constant(llvmlite.ir.IntType(32), 4))

        # Create a LLVM module object from the IR
        mod = llvm.parse_assembly(str(module))
        mod.verify()

        return mod
