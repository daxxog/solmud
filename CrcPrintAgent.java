import org.objectweb.asm.*;
import org.objectweb.asm.commons.AdviceAdapter;

import java.lang.instrument.ClassFileTransformer;
import java.lang.instrument.IllegalClassFormatException;
import java.lang.instrument.Instrumentation;
import java.security.ProtectionDomain;

public class CrcPrintAgent {
    public static void premain(String agentArgs, Instrumentation inst) {
        inst.addTransformer(new CrcTransformer());
    }

    static class CrcTransformer implements ClassFileTransformer {
        @Override
        public byte[] transform(ClassLoader loader, String className, Class<?> classBeingRedefined,
                                ProtectionDomain protectionDomain, byte[] classfileBuffer) throws IllegalClassFormatException {
            if (className.replace('/', '.').equals("client")) {
                ClassReader cr = new ClassReader(classfileBuffer);
                ClassWriter cw = new ClassWriter(ClassWriter.COMPUTE_FRAMES);
                ClassVisitor cv = new ClassVisitor(Opcodes.ASM9, cw) {
                    @Override
                    public MethodVisitor visitMethod(int access, String name, String desc, String signature, String[] exceptions) {
                        MethodVisitor mv = super.visitMethod(access, name, desc, signature, exceptions);
                        if (name.equals("h") && desc.equals("(I)V")) {
                            return new MethodVisitor(Opcodes.ASM9, mv) {
                                @Override
                                public void visitMethodInsn(int opcode, String owner, String mName, String mDesc, boolean itf) {
                                    super.visitMethodInsn(opcode, owner, mName, mDesc, itf);
                                    if (opcode == Opcodes.INVOKEVIRTUAL && mName.equals("readFully") && mDesc.equals("([BII)V") && owner.equals("java/io/DataInputStream")) {
                                        // Insert prints after readFully
                                        super.visitFieldInsn(Opcodes.GETSTATIC, "java/lang/System", "out", "Ljava/io/PrintStream;");
                                        super.visitLdcInsn("=== RECEIVED 40-BYTE CRC TABLE RAW BYTES ===");
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "println", "(Ljava/lang/String;)V", false);

                                        // Print the 40 bytes as hex from MBMGIXGO.y array (aload 6)
                                        super.visitFieldInsn(Opcodes.GETSTATIC, "java/lang/System", "out", "Ljava/io/PrintStream;");
                                        super.visitLdcInsn("Raw bytes: ");
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "print", "(Ljava/lang/String;)V", false);

                                        for (int i = 0; i < 40; i++) {
                                            super.visitFieldInsn(Opcodes.GETSTATIC, "java/lang/System", "out", "Ljava/io/PrintStream;");
                                            super.visitVarInsn(Opcodes.ALOAD, 6); // aload 6 = MBMGIXGO
                                            super.visitFieldInsn(Opcodes.GETFIELD, "MBMGIXGO", "y", "[B");
                                            super.visitLdcInsn(i);
                                            super.visitInsn(Opcodes.BALOAD);
                                            super.visitMethodInsn(Opcodes.INVOKESTATIC, "java/lang/Integer", "toHexString", "(I)Ljava/lang/String;", false);
                                            super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "print", "(Ljava/lang/String;)V", false);
                                            super.visitLdcInsn(" ");
                                            super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "print", "(Ljava/lang/String;)V", false);
                                        }

                                        super.visitFieldInsn(Opcodes.GETSTATIC, "java/lang/System", "out", "Ljava/io/PrintStream;");
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "println", "()V", false);
                                    }
                                }

                                @Override
                                public void visitInsn(int opcode) {
                                    if (opcode == Opcodes.IFEQ) { // After the hash comparison if_icmpeq L12, print if mismatch
                                        super.visitInsn(opcode);
                                        super.visitFieldInsn(Opcodes.GETSTATIC, "java/lang/System", "out", "Ljava/io/PrintStream;");
                                        super.visitLdcInsn("CRC hash mismatch! Received: ");
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "print", "(Ljava/lang/String;)V", false);
                                        super.visitVarInsn(Opcodes.ILOAD, 8); // iload 8 received hash
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "print", "(I)V", false);
                                        super.visitLdcInsn(" Computed: ");
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "print", "(Ljava/lang/String;)V", false);
                                        super.visitVarInsn(Opcodes.ILOAD, 9); // iload 9 computed hash
                                        super.visitMethodInsn(Opcodes.INVOKEVIRTUAL, "java/io/PrintStream", "println", "(I)V", false);
                                    } else {
                                        super.visitInsn(opcode);
                                    }
                                }
                            };
                        }
                        return mv;
                    }
                };
                cr.accept(cv, 0);
                return cw.toByteArray();
            }
            return classfileBuffer;
        }
    }
}
