import React from 'react';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faQuestionCircle from '@fortawesome/fontawesome-free-solid/faQuestionCircle';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import './PPMPaymentRequest.css';

const PPMPaymentRequestIntro = props => {
  const { history, match } = props;
  return (
    <div className="usa-grid ppm-payment-req-intro">
      <h3 className="title">Request PPM Payment</h3>
      <p>You'll need the following documents</p>
      <ul>
        <li>
          <strong>Weight tickets</strong> both empty & full, for each vehicle and trip{' '}
          <Link className="weight-ticket-examples-link" to="/weight-ticket-examples">
            <FontAwesomeIcon aria-hidden className="color_blue_link" icon={faQuestionCircle} />
          </Link>
        </li>
        <li>
          <strong>Storage and moving expenses</strong> (optional), such as:
          <ul>
            <li>storage</li>
            <li>tolls & weighing fees</li>
            <li>rental equipment</li>
          </ul>
        </li>
      </ul>
      <p>
        <Link to="/allowable-expenses">List of allowable expenses</Link>
      </p>
      {/* TODO: change onclick handler to go to next page in flow */}
      <PPMPaymentRequestActionBtns
        onClick={() => {
          history.push(`/moves/${match.params.moveId}/ppm-weight-ticket`);
        }}
        nextBtnLabel="Get Started"
      />
    </div>
  );
};
export default PPMPaymentRequestIntro;
