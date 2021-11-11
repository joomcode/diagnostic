package report

import (
	"archive/zip"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"runtime"
	"sort"
	"sync"
	"time"
)

//go:embed report.gohtml
var reportTemplate string

type Attachment []byte

type Report struct {
	Created   time.Time
	OSVersion string

	lock        sync.Mutex
	attachments map[string]Attachment
}

func New() *Report {
	osVersion, _ := GetSystemVersion()
	if osVersion == "" {
		osVersion = runtime.GOOS + ", " + runtime.GOARCH
	}
	return &Report{
		Created:   time.Now().Local(),
		OSVersion: osVersion,
	}
}

func (r *Report) CreateAttachment(name string, content []byte) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.attachments == nil {
		r.attachments = make(map[string]Attachment)
	}
	r.attachments[name] = content
}

func (r *Report) GenerateHtml() ([]byte, error) {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't parse report.gohtml: %v", err))
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r); err != nil {
		return nil, errors.New(fmt.Sprintf("can't render report: %v", err))
	}
	return buf.Bytes(), nil
}

func (r *Report) generateAttachments() (map[string]Attachment, error) {
	attachments := make(map[string]Attachment)
	if report, err := r.GenerateHtml(); err != nil {
		return nil, err
	} else {
		attachments["report.html"] = report
		fmt.Println(string(report))
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if r.attachments != nil {
		for k, v := range r.attachments {
			attachments[k] = v
		}
	}
	return attachments, nil
}

func (r *Report) SaveReport(ctx context.Context, output io.Writer) error {
	attachments, err := r.generateAttachments()
	if err != nil {
		return err
	}
	names := make([]string, 0, len(attachments))
	for name := range attachments {
		names = append(names, name)
	}
	sort.Strings(names)

	z := zip.NewWriter(output)
	if err := z.Flush(); err != nil {
		return err
	}
	for _, name := range names {
		header, err := z.CreateHeader(&zip.FileHeader{
			Name:     name,
			Modified: r.Created,
			Method:   zip.Deflate,
		})
		if err != nil {
			return err
		}
		if _, err := header.Write(attachments[name]); err != nil {
			return err
		}
	}
	if err := z.Close(); err != nil {
		return err
	}
	return nil
}
