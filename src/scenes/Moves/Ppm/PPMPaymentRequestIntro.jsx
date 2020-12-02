import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import './PPMPaymentRequest.css';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import { withContext } from 'shared/AppContext';
import { connect } from 'react-redux';
import Alert from 'shared/Alert';
import { reduxForm } from 'redux-form';
import styles from './PPMPaymentRequestIntro.module.scss';
import { loadPPMs, updatePPM, selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import { bindActionCreators } from 'redux';

class PPMPaymentRequestIntro extends Component {
  state = {
    ppmUpdateError: false,
  };

  componentDidMount() {
    this.props.loadPPMs(this.props.moveID);
  }

  updatePpmDate = (formValues) => {
    const { history, moveID, currentPPM } = this.props;
    if (formValues.actual_move_date && currentPPM) {
      const updatedPPM = { ...currentPPM, actual_move_date: formValues.actual_move_date };
      this.props
        .updatePPM(moveID, updatedPPM.id, updatedPPM)
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
                  <FontAwesomeIcon aria-hidden className="color_blue_link" icon="question-circle" />
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
  const currentPPM = selectActivePPMForMove(state, moveID);
  const actualMoveDate = currentPPM.actual_move_date ? currentPPM.actual_move_date : null;
  return {
    moveID: moveID,
    currentPPM: currentPPM,
    schema: get(state, 'swaggerInternal.spec.definitions.PatchPersonallyProcuredMovePayload'),
    initialValues: { actual_move_date: actualMoveDate },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      loadPPMs,
      updatePPM,
    },
    dispatch,
  );
}

export default withContext(connect(mapStateToProps, mapDispatchToProps)(PPMPaymentRequestIntro));
