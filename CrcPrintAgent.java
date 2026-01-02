import java.lang.instrument.ClassFileTransformer;
import java.lang.instrument.Instrumentation;
import java.security.ProtectionDomain;

import org.objectweb.asm.*;

public class CrcPrintAgent {

    // counter lives in agent class (always visible)
    private static int index = 0;

    public static void premain(String agentArgs, Instrumentation inst) {
        inst.addTransformer(new CrcTransformer());
        System.out.println("[CRC AGENT] Loaded");
    }

    static class CrcTransformer implements ClassFileTransformer {

        @Override
        public byte[] transform(
                ClassLoader loader,
                String className,
                Class<?> classBeingRedefined,
                ProtectionDomain protectionDomain,
                byte[] classfileBuffer) {

            if (!"MBMGIXGO".equals(className)) {
                return null;
            }

            System.out.println("[CRC AGENT] Transforming MBMGIXGO");

            ClassReader cr = new ClassReader(classfileBuffer);
            ClassWriter cw = new ClassWriter(cr, ClassWriter.COMPUTE_MAXS);

            ClassVisitor cv = new ClassVisitor(Opcodes.ASM9, cw) {
                @Override
                public MethodVisitor visitMethod(
                        int access,
                        String name,
                        String descriptor,
                        String signature,
                        String[] exceptions) {

                    MethodVisitor mv = super.visitMethod(
                            access, name, descriptor, signature, exceptions);

                    if (name.equals("h") && descriptor.equals("()I")) {
                        System.out.println("[CRC AGENT] Hooking MBMGIXGO.h()");

                        return new MethodVisitor(Opcodes.ASM9, mv) {
                            @Override
                            public void visitInsn(int opcode) {
                                if (opcode == Opcodes.IRETURN) {
                                    // duplicate return value
                                    super.visitInsn(Opcodes.DUP);

                                    // call CrcPrintAgent.log(int)
                                    super.visitMethodInsn(
                                            Opcodes.INVOKESTATIC,
                                            "CrcPrintAgent",
                                            "log",
                                            "(I)V",
                                            false
                                    );
                                }
                                super.visitInsn(opcode);
                            }
                        };
                    }
                    return mv;
                }
            };

            cr.accept(cv, 0);
            return cw.toByteArray();
        }
    }

    // This method IS visible (agent class)
    public static void log(int value) {
        if (index < 9) {
            System.out.printf(
                "EXPECTED CRC[%d] = 0x%08X%n",
                index,
                value
            );
        } else {
            System.out.printf(
                "EXPECTED CRC HASH = 0x%08X%n",
                value
            );
        }
        index++;
    }
}

