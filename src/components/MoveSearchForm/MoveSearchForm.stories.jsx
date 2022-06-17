import React from 'react';

import MoveSearchForm from './MoveSearchForm';

export default {
  title: 'Office Components/Move Search Form',
  component: MoveSearchForm,
};

export const Basic = () => {
  return (
    <div className="officeApp">
      <MoveSearchForm onSubmit={() => {}} />
    </div>
  );
};
