import React from 'react';

import PortTable from './PortTable';

const poeLocationSet = {
  poeLocation: {
    portCode: 'PDX',
    portName: 'PORTLAND INTL',
    city: 'PORTLAND',
    state: 'OREGON',
    zip: '97220',
  },
  podLocation: null,
};

const podLocationSet = {
  poeLocation: null,
  podLocation: {
    portCode: 'SEA',
    portName: 'SEATTLE TACOMA INTL',
    city: 'SEATTLE',
    state: 'WASHINGTON',
    zip: '98158',
  },
};

export default {
  title: 'Office Components / PortTable',
  component: PortTable,
};

export const standard = () => {
  return (
    <div className="officeApp">
      <PortTable {...poeLocationSet} />
    </div>
  );
};

export const podLocationDisplay = () => {
  return (
    <div className="officeApp">
      <PortTable {...podLocationSet} />
    </div>
  );
};
