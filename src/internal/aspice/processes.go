package aspice

import (
	"github.com/paulefl/req42-tracer/src/internal/model"
)

// ProcessRegistry holds all ASPICE PAM 4.0 process definitions.
type ProcessRegistry struct {
	processes map[string]*model.ASPICEProcessLevel
	bestPractices map[string]*model.ASPICEBestPractice
}

// NewProcessRegistry initializes the ASPICE PAM 4.0 process registry.
func NewProcessRegistry() *ProcessRegistry {
	pr := &ProcessRegistry{
		processes: make(map[string]*model.ASPICEProcessLevel),
		bestPractices: make(map[string]*model.ASPICEBestPractice),
	}

	pr.registerProcesses()
	pr.registerBestPractices()

	return pr
}

// registerProcesses registers all ASPICE processes.
func (pr *ProcessRegistry) registerProcesses() {
	processes := []model.ASPICEProcessLevel{
		{ID: "SYS.1", Name: "Stakeholder Requirements Definition", Description: "Define stakeholder requirements"},
		{ID: "SYS.2", Name: "System Requirements Analysis", Description: "Analyze and specify system requirements"},
		{ID: "SYS.3", Name: "System Architecture Design", Description: "Define system architecture"},

		{ID: "SWE.1", Name: "Software Requirement Analysis", Description: "Analyze software requirements derived from system requirements"},
		{ID: "SWE.2", Name: "Software Design", Description: "Design software architecture and components"},
		{ID: "SWE.3", Name: "Software Unit Implementation and Verification", Description: "Implement and verify software units"},
		{ID: "SWE.4", Name: "Software Integration and Integration Testing", Description: "Integrate and test software components"},
		{ID: "SWE.5", Name: "Software Testing", Description: "Test integrated software"},
		{ID: "SWE.6", Name: "Software Configuration Management", Description: "Manage software configuration"},
	}

	for i := range processes {
		pr.processes[processes[i].ID] = &processes[i]
	}
}

// registerBestPractices registers best practices for key processes.
func (pr *ProcessRegistry) registerBestPractices() {
	bestPractices := []model.ASPICEBestPractice{
		// SWE.1 - Software Requirement Analysis
		{ID: "SWE.1.BP1", ProcessID: "SWE.1", Title: "Develop Requirements", Description: "Software requirements are developed from stakeholder requirements"},
		{ID: "SWE.1.BP2", ProcessID: "SWE.1", Title: "Establish Traceability", Description: "Establish bidirectional traceability to system requirements"},
		{ID: "SWE.1.BP3", ProcessID: "SWE.1", Title: "Analyze Requirements", Description: "Analyze software requirements for correctness and completeness"},
		{ID: "SWE.1.BP4", ProcessID: "SWE.1", Title: "Manage Consistency", Description: "Manage consistency between software requirements and architecture"},
		{ID: "SWE.1.BP5", ProcessID: "SWE.1", Title: "Ensure Consistency", Description: "Ensure consistency with stakeholder and system requirements"},
		{ID: "SWE.1.BP6", ProcessID: "SWE.1", Title: "Ensure Testability", Description: "Ensure software requirements are testable"},
		{ID: "SWE.1.BP7", ProcessID: "SWE.1", Title: "Review Requirements", Description: "Conduct reviews of software requirements"},
		{ID: "SWE.1.BP8", ProcessID: "SWE.1", Title: "Manage Bidirectional Traceability", Description: "Manage bidirectional traceability between requirements and design"},

		// SWE.2 - Software Design
		{ID: "SWE.2.BP1", ProcessID: "SWE.2", Title: "Develop Design", Description: "Design is developed from software requirements"},
		{ID: "SWE.2.BP2", ProcessID: "SWE.2", Title: "Establish Traceability", Description: "Establish bidirectional traceability to software requirements"},
		{ID: "SWE.2.BP3", ProcessID: "SWE.2", Title: "Analyze Design", Description: "Analyze design for correctness and completeness"},
		{ID: "SWE.2.BP4", ProcessID: "SWE.2", Title: "Ensure Traceability", Description: "Ensure traceability to software requirements"},
		{ID: "SWE.2.BP5", ProcessID: "SWE.2", Title: "Review Design", Description: "Conduct reviews of software design"},

		// SWE.3 - Software Unit Implementation and Verification
		{ID: "SWE.3.BP1", ProcessID: "SWE.3", Title: "Implement Units", Description: "Software units are implemented from design"},
		{ID: "SWE.3.BP2", ProcessID: "SWE.3", Title: "Verify Units", Description: "Software units are verified against design"},
		{ID: "SWE.3.BP3", ProcessID: "SWE.3", Title: "Establish Traceability", Description: "Establish traceability between implementation and design"},
		{ID: "SWE.3.BP4", ProcessID: "SWE.3", Title: "Manage Traceability", Description: "Manage traceability throughout software unit implementation"},

		// SWE.5 - Software Testing
		{ID: "SWE.5.BP1", ProcessID: "SWE.5", Title: "Plan Testing", Description: "Test strategy and plan are established"},
		{ID: "SWE.5.BP2", ProcessID: "SWE.5", Title: "Implement Tests", Description: "Test cases and procedures are implemented"},
		{ID: "SWE.5.BP3", ProcessID: "SWE.5", Title: "Establish Traceability", Description: "Establish traceability between tests and requirements"},
		{ID: "SWE.5.BP4", ProcessID: "SWE.5", Title: "Execute Tests", Description: "Test cases are executed"},
		{ID: "SWE.5.BP5", ProcessID: "SWE.5", Title: "Track Defects", Description: "Defects are identified, tracked, and resolved"},
	}

	for i := range bestPractices {
		pr.bestPractices[bestPractices[i].ID] = &bestPractices[i]
	}
}

// GetProcess returns a process definition by ID.
func (pr *ProcessRegistry) GetProcess(id string) *model.ASPICEProcessLevel {
	return pr.processes[id]
}

// GetBestPractice returns a best practice definition by ID.
func (pr *ProcessRegistry) GetBestPractice(id string) *model.ASPICEBestPractice {
	return pr.bestPractices[id]
}

// ListProcesses returns all registered processes.
func (pr *ProcessRegistry) ListProcesses() []*model.ASPICEProcessLevel {
	var processes []*model.ASPICEProcessLevel
	for _, p := range pr.processes {
		processes = append(processes, p)
	}
	return processes
}

// ListBestPracticesForProcess returns all best practices for a given process.
func (pr *ProcessRegistry) ListBestPracticesForProcess(processID string) []*model.ASPICEBestPractice {
	var bps []*model.ASPICEBestPractice
	for _, bp := range pr.bestPractices {
		if bp.ProcessID == processID {
			bps = append(bps, bp)
		}
	}
	return bps
}
