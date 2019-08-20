import React from 'react';
import weightTixExample from 'shared/images/weight_tix_example.png';
import weightScenario1 from 'shared/images/weight_scenario1.png';
import weightScenario2 from 'shared/images/weight_scenario2.png';
import weightScenario3 from 'shared/images/weight_scenario3.png';
import weightScenario4 from 'shared/images/weight_scenario4.png';
import weightScenario5 from 'shared/images/weight_scenario5.png';

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
        <div className="subheader">
          You need two weight tickets for each trip you took: one with the vehicle empty, one with it full.
        </div>
        <img className="weight-ticket-example-img" alt="weight ticket example" src={weightTixExample} /> = A{' '}
        <strong>trip</strong> includes both an empty and <strong>full</strong> weight ticket
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 1</div>
        <p>You and your spouse each drove a vehicle filled with stuff to your destination.</p>
        <div className="usa-width-one-whole">
          <img src={weightScenario1} alt="weight scenario 1" />
        </div>
        <p>
          You must upload weight tickets for <strong>2 trips</strong>, which is <strong>4 tickets</strong> total.
        </p>
        <p className="text-gray secondary-label">
          <em>That's one empty weight ticket and one full weight ticket for each vehicle.</em>
        </p>
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 2</div>
        <p>You made two separate trips in one vehicle to bring stuff to your destination.</p>
        <div className="usa-width-one-whole">
          <img src={weightScenario2} alt="weight scenario 2" />
        </div>
        <p>
          You must upload weight tickets for <strong>2 trips</strong>, which is <strong>4 tickets</strong> total.
        </p>
        <p className="secondary-label">
          <em>
            <span className="text-gray">
              That's one empty and one full ticket for the first trip, and one empty and one full for the second trip.
            </span>{' '}
            You do need to weigh your empty vehicle a second time.
          </em>
        </p>
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 3</div>
        <p>
          You and your spouse each drove a vehicle filled with stuff to your destination. Then you made a second trip in
          your vehicle (without your spouse) to bring more stuff.
        </p>
        <div className="usa-width-one-whole">
          <img src={weightScenario3} alt="weight scenario 3" />
        </div>
        <p>
          You must upload weight tickets for <strong>3 trips</strong>, which is <strong>6 tickets</strong> total.
        </p>
        <p className="text-gray secondary-label">
          <em>
            That's one empty and one full weight ticket for each vehicle on the first trip, and an empty and full weight
            ticket for your vehicle on the second trip.
          </em>
        </p>
      </section>
      <div className="dashed-divider" />
      <section>
        <div className="subheader">Scenario 4</div>
        <p>
          You drove your car with an attached rental trailer to your destination, then made a second trip to bring more
          stuff.
        </p>
        <div className="usa-width-one-whole">
          <img src={weightScenario4} alt="weight scenario 4" />
        </div>
        <p>
          You must upload weight tickets for <strong>2 trips</strong>, which is <strong>4 tickets</strong> total.
        </p>
        <p className="text-gray secondary-label">
          <em>
            You can't claim the weight of your rented trailer in your move. All weight tickets must include the weight
            of your car with trailer attached.
          </em>
        </p>
      </section>
      <div className="dashed-divider" />

      <section>
        <div className="subheader">Scenario 5</div>
        <p>
          You drove your car with an attached trailer that you own, and that meets the trailer criteria, to make two
          separate trips to move stuff to your destination.
        </p>
        <div className="usa-width-one-whole">
          <img src={weightScenario5} alt="weight scenario 5" />
        </div>
        <p>
          You must upload weight tickets for <strong>2 trips</strong>, which is <strong>4 tickets</strong> total.
        </p>
        <p className="text-gray secondary-label">
          <em>
            You can claim the weight of your own trailer once per move (not per trip). The empty weight ticket for your
            first trip should be the weight of your car only. All 3 additional weight tickets should include the weight
            of your car with trailer attached.
          </em>
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
