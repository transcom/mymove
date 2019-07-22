import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faQuestionCircle from '@fortawesome/fontawesome-free-solid/faQuestionCircle';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import './PPMPaymentRequest.css';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import { withContext } from 'shared/AppContext';
import { connect } from 'react-redux';
import { updatePPM } from 'shared/Entities/modules/ppms';
import Alert from 'shared/Alert';
import { reduxForm } from 'redux-form';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';
import './PPMPaymentRequestIntro.css';

class PPMPaymentRequestIntro extends Component {
  state = {
    ppmUpdateError: false,
  };

  updatePpmDate = formValues => {
    const {
      history,
      moveId,
      currentPpm: { id: ppmId },
    } = this.props;
    if (formValues.actual_move_date) {
      this.props
        .updatePPM(moveId, ppmId, { actual_move_date: formValues.actual_move_date })
        .then(() => history.push(`/moves/${moveId}/ppm-weight-ticket`))
        .catch(() => {
          this.setState({ ppmUpdateError: true });
        });
    }
  };

  render() {
    const { schema, invalid, handleSubmit, submitting } = this.props;
    const { ppmUpdateError } = this.state;
    return (
      <div className="usa-grid ppm-payment-req-intro">
        {ppmUpdateError && (
          <div className="usa-grid">
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                Something went wrong contacting the server.
              </Alert>
            </div>
          </div>
        )}
        <h3 className="title">Request PPM Payment</h3>
        <p>Gather these documents:</p>
        <ul>
          <li>
            <strong>Weight tickets,</strong> both empty & full, for <em>each</em> vehicle and trip{' '}
            <Link className="weight-ticket-examples-link" to="/weight-ticket-examples">
              <FontAwesomeIcon aria-hidden className="color_blue_link" icon={faQuestionCircle} />
            </Link>
          </li>
          <li>
            <strong>Storage and moving expenses</strong> (if used), such as:
            <ul>
              <li>storage</li>
              <li>tolls & weighing fees</li>
              <li>rental equipment</li>
            </ul>
          </li>
        </ul>
        <p>
          <Link to="/allowable-expenses">More about expenses</Link>
        </p>
        <SwaggerField
          className="ppm-payment-request-actual-date"
          title="What day did you depart?"
          fieldName="actual_move_date"
          swagger={schema}
          required
        />
        <PPMPaymentRequestActionBtns
          hasConfirmation={true}
          saveAndAddHandler={handleSubmit(this.updatePpmDate)}
          submitButtonsAreDisabled={invalid || submitting}
          nextBtnLabel="Get Started"
        />
      </div>
    );
  }
}

const formName = 'ppm_payment_intro_wizard';
PPMPaymentRequestIntro = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PPMPaymentRequestIntro);

function mapStateToProps(state, ownProps) {
  const moveId = ownProps.match.params.moveId;
  const currentPpm = selectPPMForMove(state, moveId) || get(state, 'ppm.currentPpm');
  const actualMoveDate = currentPpm.actual_move_date ? currentPpm.actual_move_date : null;
  console.log(actualMoveDate);
  return {
    moveId: moveId,
    currentPpm: get(state, 'ppm.currentPpm'),
    schema: get(state, 'swaggerInternal.spec.definitions.PatchPersonallyProcuredMovePayload'),
    initialValues: { actual_move_date: actualMoveDate },
  };
}

const mapDispatchToProps = {
  updatePPM,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(PPMPaymentRequestIntro));
