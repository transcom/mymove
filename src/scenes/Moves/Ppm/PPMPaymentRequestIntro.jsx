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
import Alert from 'shared/Alert';
import { reduxForm } from 'redux-form';
import styles from './PPMPaymentRequestIntro.module.scss';
import { createOrUpdatePpm } from './ducks';

class PPMPaymentRequestIntro extends Component {
  state = {
    ppmUpdateError: false,
  };

  updatePpmDate = formValues => {
    const { history, moveId, currentPpm } = this.props;
    if (formValues.actual_move_date && currentPpm) {
      const updatedPpm = { ...currentPpm, actual_move_date: formValues.actual_move_date };
      this.props
        .createOrUpdatePpm(moveId, updatedPpm)
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
      <div className="grid-container usa-prose ppm-payment-req-intro">
        {ppmUpdateError && (
          <div className="grid-row">
            <div className="grid-col-12 error-message">
              <Alert type="error" heading="An error occurred">
                Something went wrong contacting the server.
              </Alert>
            </div>
          </div>
        )}
        <div className="grid-row">
          <div className="grid-col-12">
            <h3 className="title">Request PPM Payment</h3>
            <p>Gather these documents:</p>
            <ul>
              <li>
                <strong>Weight tickets,</strong> both empty & full, for <em>each</em> vehicle and trip{' '}
                <Link className="weight-ticket-examples-link usa-link" to="/weight-ticket-examples">
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
              <Link to="/allowable-expenses" className="usa-link">
                More about expenses
              </Link>
            </p>
            <SwaggerField
              className={styles['ppm-payment-request-actual-date']}
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
        </div>
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
  const currentPpm = get(state, 'ppm.currentPpm');
  const actualMoveDate = currentPpm.actual_move_date ? currentPpm.actual_move_date : null;
  return {
    moveId: moveId,
    currentPpm: currentPpm,
    schema: get(state, 'swaggerInternal.spec.definitions.PatchPersonallyProcuredMovePayload'),
    initialValues: { actual_move_date: actualMoveDate },
  };
}

const mapDispatchToProps = {
  createOrUpdatePpm,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(PPMPaymentRequestIntro));
