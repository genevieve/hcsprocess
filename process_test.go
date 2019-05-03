package hcsprocess_test

import (
	"github.com/genevieve/hcsprocess"
	"github.com/genevieve/hcsprocess/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Process", func() {
	var wrappedProcess *hcsprocess.Process
	var fakeProcess *fakes.Process

	BeforeEach(func() {
		fakeProcess = &fakes.Process{}
		wrappedProcess = hcsprocess.New(fakeProcess)
	})

	Describe("AttachIO", func() {
		var (
			attachedStdin,
			attachedStdout,
			attachedStderr,
			processStdin,
			processStdout,
			processStderr *gbytes.Buffer
		)

		BeforeEach(func() {
			attachedStdin = gbytes.BufferWithBytes([]byte("something-on-stdin"))
			attachedStdout = gbytes.NewBuffer()
			attachedStderr = gbytes.NewBuffer()

			processStdin = gbytes.NewBuffer()
			processStdout = gbytes.BufferWithBytes([]byte("something-on-stdout"))
			processStderr = gbytes.BufferWithBytes([]byte("something-on-stderr"))

			fakeProcess.StdioCall.Returns.Stdin = processStdin
			fakeProcess.StdioCall.Returns.Stdout = processStdout
			fakeProcess.StdioCall.Returns.Stderr = processStderr
		})

		It("attaches process IO to stdin, stdout, and stderr", func() {
			exitCode, err := wrappedProcess.AttachIO(attachedStdin, attachedStdout, attachedStderr)
			Expect(err).NotTo(HaveOccurred())

			Expect(exitCode).To(Equal(0))

			Eventually(processStdin).Should(gbytes.Say("something-on-stdin"))
			Eventually(attachedStdout).Should(gbytes.Say("something-on-stdout"))
			Eventually(attachedStderr).Should(gbytes.Say("something-on-stderr"))
		})

		It("closes the process' stdin pipe after copying", func() {
			exitCode, err := wrappedProcess.AttachIO(attachedStdin, nil, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(0))

			Expect(fakeProcess.CloseStdinCall.CallCount).To(Equal(1))

			Eventually(processStdin).Should(gbytes.Say("something-on-stdin"))
		})

		Context("when io.Copy is blocking", func() {
			BeforeEach(func() {
				fakeProcess.StdioCall.Returns.Stdin = processStdin
			})
			It("exits after 5 seconds", func() {
				exitCode, err := wrappedProcess.AttachIO(attachedStdin, nil, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(exitCode).To(Equal(0))

				Expect(fakeProcess.CloseStdinCall.CallCount).To(Equal(1))
			})
		})
	})
})
