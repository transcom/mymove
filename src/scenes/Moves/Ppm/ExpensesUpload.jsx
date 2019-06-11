import React, { Component } from 'react';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';

import 'shared/DocumentViewer/DocumentUploader.jsx';
import { get, isEmpty } from 'lodash';
import RadioButton from 'shared/RadioButton';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faQuestionCircle from '@fortawesome/fontawesome-free-solid/faQuestionCircle';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Uploader from 'shared/Uploader';
import Checkbox from 'shared/Checkbox';
import { getFormValues, reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import Alert from 'shared/Alert';

class ExpensesUpload extends Component {
  state = { ...this.initialState, expenseNumber: 1 };

  static nextBtnLabels = {
    SaveAndAddAnother: 'Save & Add Another',
    SaveAndContinue: 'Save & Continue',
  };

  static paymentMethods = {
    Other: 'OTHER',
    GTCC: 'GTCC',
  };

  get initialState() {
    return {
      paymentMethod: ExpensesUpload.paymentMethods.Other,
      uploaderIsIdle: true,
      missingReceipt: false,
      expenseType: '',
      haveMoreExpenses: 'No',
      moveDocumentCreateError: false,
    };
  }

  handleHasMoreExpensesChange = event => {
    this.setState({
      haveMoreExpenses: event.target.value,
    });
  };

  handlePaymentMethodChange = event => {
    this.setState({
      paymentMethod: event.target.value,
    });
  };

  onAddFile = () => {
    this.setState({
      uploaderIsIdle: false,
    });
  };

  onChange = (newUploads, uploaderIsIdle) => {
    this.setState({
      uploaderIsIdle,
    });
  };

  handleCheckboxChange = event => {
    this.setState({
      [event.target.name]: event.target.checked,
    });
  };

  handleHowDidYouPayForThis = () => {
    alert('Cash, personal credit card, check — any payment method that’s not your GTCC.');
  };

  render() {
    const { missingReceipt, paymentMethod, haveMoreExpenses, expenseNumber, moveDocumentCreateError } = this.state;
    const { moveDocSchema, formValues, isPublic } = this.props;
    const nextBtnLabel =
      haveMoreExpenses === 'Yes'
        ? ExpensesUpload.nextBtnLabels.SaveAndAddAnother
        : ExpensesUpload.nextBtnLabels.SaveAndContinue;
    const hasMovingExpenseType = !isEmpty(formValues) && formValues.moving_expense_type !== '';
    return (
      <>
        <WizardHeader
          title="Expenses"
          right={
            <ProgressTimeline>
              <ProgressTimelineStep name="Weight" completed />
              <ProgressTimelineStep name="Expenses" current />
              <ProgressTimelineStep name="Review" />
            </ProgressTimeline>
          }
        />

        <div className="usa-grid expenses-container">
          <h3 className="expenses-header">Expense {expenseNumber}</h3>
          <p>
            Upload expenses one at a time.{' '}
            <Link to="/allowable-expenses">
              <FontAwesomeIcon aria-hidden className="color_blue_link" icon={faQuestionCircle} />
            </Link>
          </p>
          <form>
            {moveDocumentCreateError && (
              <div className="usa-grid">
                <div className="usa-width-one-whole error-message">
                  <Alert type="error" heading="An error occurred">
                    Something went wrong contacting the server.
                  </Alert>
                </div>
              </div>
            )}
            <SwaggerField title="Expense type" fieldName="moving_expense_type" swagger={moveDocSchema} required />
            {hasMovingExpenseType && (
              <>
                <SwaggerField title="Document title" fieldName="title" swagger={moveDocSchema} required />
                <SwaggerField title="Amount" fieldName="requested_amount_cents" swagger={moveDocSchema} required />
                <div className="usa-width-two-thirds uploader-wrapper">
                  <Uploader
                    isPublic={isPublic}
                    onRef={ref => (this.uploader = ref)}
                    onChange={this.onChange}
                    onAddFile={this.onAddFile}
                  />
                </div>
                <Checkbox
                  label="I'm missing this receipt"
                  name="missingReceipt"
                  checked={missingReceipt}
                  onChange={this.handleCheckboxChange}
                  normalizeLabel
                />
                <div className="radio-group-wrapper" style={{ marginBottom: '0px' }}>
                  <p className="radio-group-header">How did you pay for this?</p>
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Government travel charge card (GTCC)"
                    value={ExpensesUpload.paymentMethods.GTCC}
                    name="payment_method"
                    checked={paymentMethod === ExpensesUpload.paymentMethods.GTCC}
                    onChange={this.handlePaymentMethodChange}
                  />
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Other"
                    value={ExpensesUpload.paymentMethods.Other}
                    name="payment_method"
                    checked={paymentMethod === ExpensesUpload.paymentMethods.Other}
                    onChange={this.handlePaymentMethodChange}
                  />
                  <FontAwesomeIcon
                    aria-hidden
                    className="color_blue_link"
                    icon={faQuestionCircle}
                    onClick={this.handleHowDidYouPayForThis}
                  />
                </div>
                <div className="dashed-divider" style={{ marginTop: '0px' }} />
                <div className="radio-group-wrapper">
                  <p className="radio-group-header">Do you have more expenses to upload?</p>
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Yes"
                    value="Yes"
                    name="has_more_expenses"
                    checked={haveMoreExpenses === 'Yes'}
                    onChange={this.handleHasMoreExpensesChange}
                  />
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="No"
                    value="No"
                    name="has_more_expenses"
                    checked={haveMoreExpenses === 'No'}
                    onChange={this.handleHasMoreExpensesChange}
                  />
                </div>
              </>
            )}
            <PPMPaymentRequestActionBtns
              nextBtnLabel={nextBtnLabel}
              submitButtonsAreDisabled={true}
              saveForLaterHandler={() => {}}
              saveAndAddHandler={() => {}}
              displaySaveForLater={true}
            />
          </form>
        </div>
      </>
    );
  }
}

const formName = 'expense_document_upload';
ExpensesUpload = reduxForm({ form: formName })(ExpensesUpload);

function mapStateToProps(state, props) {
  const moveId = props.match.params.moveId;
  return {
    moveId: moveId,
    formValues: getFormValues(formName)(state),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    initialValues: {},
    currentPpm: get(state, 'ppm.currentPpm'),
  };
}
export default connect(mapStateToProps)(ExpensesUpload);
