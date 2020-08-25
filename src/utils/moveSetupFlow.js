// Concat service member's full name for display
export function getFullSMName(serviceMember) {
  if (!serviceMember) return '';
  return `${serviceMember.first_name} ${serviceMember.middle_name || ''} ${serviceMember.last_name} ${
    serviceMember.suffix || ''
  }`;
}

// Concat agent's full name for display
export function getFullAgentName(agent) {
  if (!agent) return '';
  return `${agent.firstName || ''} ${agent.lastName || ''}`;
}
