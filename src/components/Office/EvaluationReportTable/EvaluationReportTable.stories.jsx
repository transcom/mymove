import React from 'react';

import EvaluationReportTable from './EvaluationReportTable';

export default {
  title: 'Office Components/EvaluationReportTable',
  component: EvaluationReportTable,
};

const reports = [{ id: '12354' }];

export const empty = () => (
  <div className="officeApp">
    <EvaluationReportTable reports={[]} />
  </div>
);

export const single = () => (
  <div className="officeApp">
    <EvaluationReportTable reports={reports} />
  </div>
);
