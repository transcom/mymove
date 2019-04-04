import React from 'react';
import weightTixExample from 'shared/images/weight_tix_example.png';
import weightScenario1 from 'shared/images/weight_scenario1.png';
import weightScenario2 from 'shared/images/weight_scenario2.png';
import weightScenario3 from 'shared/images/weight_scenario3.png';

function WeightTicketExamples(props) {
  function goBack() {
    props.history.goBack();
  }
  return (
    <div className="usa-grid weight-ticket-example-container">
      <div>
        <a onClick={goBack}>{'<'} Back</a>
      </div>
      <h3 className="title">Example weight ticket scenarios</h3>
      <section>
        <div className="subheader">You need an empty and full weight ticket for each trip you took.</div>
        <img className="weight-ticket-example-img" alt="weight ticket example" src={weightTixExample} /> = A{' '}
        <strong>trip</strong> includes both an empty and <strong>full</strong> weight ticket
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 1</div>
        <p>You and your spouse each drove a vehicle filled with stuff to your destination</p>
        <div className="usa-width-one-whole">
          <img src={weightScenario1} alt="weight scenario 1" />
        </div>
        <p>
          This means you have to upload weight tickets for <strong>2 trips</strong> (or 4 tickets total).
        </p>
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 2</div>
        <p>You made two separate trips in one vehicle to bring stuff to your destination</p>
        <div className="usa-width-one-whole">
          <img src={weightScenario2} alt="weight scenario 2" />
        </div>
        <p>
          This means you have to upload weight tickets for <strong>2 trips</strong> (or 4 tickets total).
        </p>
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 3</div>
        <p>
          You and your spouse each drove a vehicle filled with stuff to your destination. Then, you made a second trip
          in your vehicle to bring more stuff.
        </p>
        <div className="usa-width-one-whole">
          <img src={weightScenario3} alt="weight scenario 3" />
        </div>
        <p>
          This means you have to upload weight tickets for <strong>3 trips</strong> (or 6 tickets total).
        </p>
      </section>
      <div className="usa-grid button-bar">
        <button className="usa-button-secondary" onClick={goBack}>
          Back
        </button>
      </div>
    </div>
  );
}

export default WeightTicketExamples;
