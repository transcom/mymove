import React from 'react';
import { useNavigate } from 'react-router-dom';

function TrailerCriteria() {
  const navigate = useNavigate();
  function goBack() {
    navigate(-1);
  }
  return (
    <div className="usa-grid trailer-criteria-container">
      <div>
        <a onClick={goBack} className="usa-link">
          {'<'} Back
        </a>
      </div>
      <h1 className="title">Trailer Criteria</h1>
      <section>
        <p>
          During your move, if you used a trailer owned by you or your spouse, you can claim its weight{' '}
          <strong>once</strong> per move if it meets these specifications:
        </p>
        <div className="list-header">
          <p> A utility trailer:</p>
          <ul>
            <li>With or without a tilt bed</li>
            <li>Single axle</li>
            <li>No more than 12 feet long from rear to trailer hitch</li>
            <li>No more than 8 feet wide from outside tire to outside tire</li>
            <li>Side rails and body no higher than 28 inches (unless detachable) </li>
            <li>Ramp or gate for the utility trailer no higher than 4 feet (unless detachable)</li>
          </ul>
        </div>
      </section>
      <p>
        You will also have to provide proof of ownership, either a registration or bill of sale. If these are
        unavailable in your state, you can provide a signed and dated statement certifying that you or your spouse own
        the trailer.
      </p>
      <div className="usa-grid button-bar">
        <button className="usa-button usa-button--secondary" onClick={goBack}>
          Back
        </button>
      </div>
    </div>
  );
}

export default TrailerCriteria;
