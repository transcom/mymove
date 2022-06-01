import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { reduxForm } from 'redux-form';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import './PPMPaymentRequest.css';
import styles from './PPMPaymentRequestIntro.module.scss';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { formatDateForSwagger } from 'shared/dates';
import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert';
import { getPPMsForMove, patchPPM } from 'services/internalApi';
import { updatePPMs, updatePPM } from 'store/entities/actions';
import { selectCurrentPPM } from 'store/entities/selectors';

class PPMPaymentRequestIntro extends Component {
  state = {
    ppmUpdateError: false,
  };

  componentDidMount() {
    getPPMsForMove(this.props.moveID).then((response) => this.props.updatePPMs(response));
  }

  updatePpmDate = (formValues) => {
    const { history, moveID, currentPPM } = this.props;
    if (formValues.actual_move_date && currentPPM) {
      const updatedPPM = { ...currentPPM, actual_move_date: formatDateForSwagger(formValues.actual_move_date) };

      patchPPM(moveID, updatedPPM)
        .then((response) => this.props.updatePPM(response))
        .then(() => history.push(`/moves/${moveID}/ppm-weight-ticket`))
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
            <h1 className="title">Request PPM Payment</h1>
            <p>Gather these documents:</p>
            <ul>
              <li>
                <strong>Weight tickets,</strong> both empty & full, for <em>each</em> vehicle and trip{' '}
                <Link className="weight-ticket-examples-link usa-link" to="/weight-ticket-examples">
                  <FontAwesomeIcon aria-hidden className="color_blue_link" icon="circle-question" />
                </Link>
              </li>
              <li>
                <strong>Storage and moving expenses</strong> (if used), such as:
                <ul>
                  <li>storage</li>
                  <li>tolls & weighing fees</li>
                  <li>rental equipment</li>
                  <li>fees for movers you hired</li>
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
  const moveID = ownProps.match.params.moveId;
  const currentPPM = selectCurrentPPM(state) || {};
  const actualMoveDate = currentPPM.actual_move_date ? currentPPM.actual_move_date : null;
  return {
    moveID,
    currentPPM,
    schema: get(state, 'swaggerInternal.spec.definitions.PatchPersonallyProcuredMovePayload'),
    initialValues: { actual_move_date: actualMoveDate },
  };
}

const mapDispatchToProps = {
  updatePPM,
  updatePPMs,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(PPMPaymentRequestIntro));
