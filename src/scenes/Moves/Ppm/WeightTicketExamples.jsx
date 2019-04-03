import React from 'react';

function WeightTicketExamples(props) {
  function goBack() {
    props.history.goBack();
  }
  return (
    <div className="usa-grid">
      <div>
        <a onClick={goBack}>{'<'} Back</a>
      </div>
      Examples here
    </div>
  );
}

export default WeightTicketExamples;
