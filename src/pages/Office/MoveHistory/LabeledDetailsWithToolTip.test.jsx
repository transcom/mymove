import React from 'react';
import { shallow } from 'enzyme';
import { render, screen } from '@testing-library/react';

import LabeledDetailsWithToolTip from './LabeledDetailsWithToolTip';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('LabeledDetailsWithToolTip', () => {
  it('renders without crashing', () => {
    const historyRecord = {
      changedValues: {
        changedValues: {
          sit_entry_date: '2023-10-01',
          shipment_id_display: '1A2B3',
          shipment_locator: 'ABC123-01',
        },
      },
    };
    const wrapper = shallow(<LabeledDetailsWithToolTip historyRecord={historyRecord} />);
    expect(wrapper.exists()).toBe(true);
  });

  it('displays field details with tooltips', () => {
    const historyRecord = {
      changedValues: {
        sit_entry_date: '2023-10-01',
        shipment_id_display: '1A2B3',
      },
      oldValues: {
        sit_entry_date: '2023-09-01',
      },
    };
    const toolTipText = 'Some tooltip text';
    const toolTipColor = 'black';
    const toolTipTextPosition = 'top';
    const toolTipIcon = 'circle-question';

    const wrapper = shallow(
      <LabeledDetailsWithToolTip
        historyRecord={historyRecord}
        toolTipText={toolTipText}
        toolTipColor={toolTipColor}
        toolTipTextPosition={toolTipTextPosition}
        toolTipIcon={toolTipIcon}
      />,
    );

    expect(wrapper.find('div')).toHaveLength(1);
    expect(wrapper.find('ToolTip').props().text).toEqual(toolTipText);
    expect(wrapper.find('ToolTip').props().color).toEqual(toolTipColor);
    expect(wrapper.find('ToolTip').props().position).toEqual(toolTipTextPosition);
    expect(wrapper.find('ToolTip').props().icon).toEqual(toolTipIcon);
  });
});

it('renders shipment_type as a header & SIT entry date', async () => {
  const historyRecord = {
    changedValues: {
      sit_entry_date: '2023-10-01',
      shipment_id_display: '1A2B3',
      shipment_locator: 'ABC123-01',
      shipment_type: SHIPMENT_OPTIONS.HHG,
    },
    oldValues: {
      sit_entry_date: '2023-09-01',
    },
  };

  render(<LabeledDetailsWithToolTip historyRecord={historyRecord} />);

  expect(screen.getByText('HHG shipment #ABC123-01')).toBeInTheDocument();
  expect(screen.getByText('SIT entry date')).toBeInTheDocument();
  expect(screen.getByText(': 01 Oct 2023')).toBeInTheDocument();
});
