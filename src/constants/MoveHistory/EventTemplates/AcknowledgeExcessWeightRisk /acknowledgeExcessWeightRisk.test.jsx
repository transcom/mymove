import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/AcknowledgeExcessWeightRisk /acknowledgeExcessWeightRisk';

describe('when given an Acknowledge excess weight risk history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'acknowledgeExcessWeightRisk',
    tableName: 'moves',
  };

  it('correctly matches the Acknowledge excess weight risk template', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });

  it('renders the proper message in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Dismissed excess weight alert')).toBeInTheDocument();
  });
});
