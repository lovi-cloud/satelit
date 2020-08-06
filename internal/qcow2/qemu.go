package qcow2

import (
	"context"
	"fmt"
	"os/exec"
)

// NOTE(whywaita): move to go-os-brick later.

func qemuimgConvertBase(ctx context.Context, args []string) ([]byte, error) {
	c := []string{"convert"}
	a := append(c, args...)
	out, err := exec.CommandContext(ctx, "qemu-img", a...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute qemu-img convert command (args: %s): %w", args, err)
	}

	return out, nil
}

// ToRaw convert os image.
func ToRaw(ctx context.Context, src, dest string) error {
	args := []string{"-O", "raw", "-t", "none", "-f", "qcow2", src, dest}

	out, err := qemuimgConvertBase(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute convert command (out: %s): %w", string(out), err)
	}

	return nil
}
