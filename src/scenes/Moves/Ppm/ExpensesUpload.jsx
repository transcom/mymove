import React, { Component } from 'react';
import { get, isEmpty, map } from 'lodash';
import { withLastLocation } from 'react-router-last-location';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { getFormValues, reduxForm } from 'redux-form';
import { connect } from 'react-redux';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import DocumentsUploaded from './PaymentReview/DocumentsUploaded';
import WizardHeader from '../WizardHeader';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { convertDollarsToCents } from 'shared/utils';
import RadioButton from 'shared/RadioButton';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import Uploader from 'shared/Uploader';
import Checkbox from 'shared/Checkbox';
import {
  createMovingExpenseDocument,
  selectPPMCloseoutDocumentsForMove,
} from 'shared/Entities/modules/movingExpenseDocuments';
import Alert from 'shared/Alert';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { withContext } from 'shared/AppContext';
import { documentSizeLimitMsg } from 'shared/constants';
import { selectCurrentPPM } from 'store/entities/selectors';

const nextPagePath = '/ppm-payment-review';
const nextBtnLabels = {
  SaveAndAddAnother: 'Save & Add Another',
  SaveAndContinue: 'Save & Continue',
};

const paymentMethods = {
  Other: 'OTHER',
  GTCC: 'GTCC',
};

const uploadReceipt =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload receipt</span></span>';

class ExpensesUpload extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      paymentMethod: paymentMethods.GTCC,
      uploaderIsIdle: true,
      missingReceipt: false,
      expenseType: '',
      haveMoreExpenses: 'No',
      moveDocumentCreateError: false,
    };
  }

  componentDidMount() {
    const { moveId } = this.props;
    this.props.getMoveDocumentsForMove(moveId);
  }

  handleRadioChange = (event) => {
    this.setState({
      [event.target.name]: event.target.value,
    });
  };

  skipHandler = () => {
    const { moveId, history } = this.props;
    history.push(`/moves/${moveId}${nextPagePath}`);
  };

  isStorageExpense = (formValues) => {
    return !isEmpty(formValues) && formValues.moving_expense_type === 'STORAGE';
  };

  saveAndAddHandler = (formValues) => {
    const { moveId, currentPpm, history } = this.props;
    const { paymentMethod, missingReceipt, haveMoreExpenses } = this.state;
    const {
      storage_start_date,
      storage_end_date,
      moving_expense_type: movingExpenseType,
      requested_amount_cents: requestedAmountCents,
    } = formValues;

    let files = this.uploader.getFiles();
    const uploadIds = map(files, 'id');
    const personallyProcuredMoveId = currentPpm ? currentPpm.id : null;
    const title = this.isStorageExpense(formValues) ? 'Storage Expense' : formValues.title;
    return this.props
      .createMovingExpenseDocument({
        moveId,
        personallyProcuredMoveId,
        uploadIds,
        title,
        movingExpenseType,
        moveDocumentType: 'EXPENSE',
        requestedAmountCents: convertDollarsToCents(requestedAmountCents),
        paymentMethod,
        notes: '',
        missingReceipt,
        storage_start_date,
        storage_end_date,
      })
      .then(() => {
        this.cleanup();
        if (haveMoreExpenses === 'No') {
          history.push(`/moves/${moveId}${nextPagePath}`);
        }
      })
      .catch((e) => {
        this.setState({ moveDocumentCreateError: true });
      });
  };

  cleanup = () => {
    const { reset } = this.props;
    this.uploader.clearFiles();
    reset();
    this.setState({ ...this.initialState });
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

  handleCheckboxChange = (event) => {
    this.setState({
      [event.target.name]: event.target.checked,
    });
  };

  handleHowDidYouPayForThis = () => {
    alert('Cash, personal credit card, check — any payment method that’s not your GTCC.');
  };

  isInvalidUploaderState = () => {
    const { missingReceipt } = this.state;
    const receiptUploaded = this.uploader && !this.uploader.isEmpty();
    return missingReceipt === receiptUploaded;
  };

  render() {
    const { missingReceipt, paymentMethod, haveMoreExpenses, moveDocumentCreateError } = this.state;
    const { moveDocSchema, formValues, isPublic, handleSubmit, submitting, expenses, expenseSchema, invalid, moveId } =
      this.props;
    const nextBtnLabel = haveMoreExpenses === 'Yes' ? nextBtnLabels.SaveAndAddAnother : nextBtnLabels.SaveAndContinue;
    const hasMovingExpenseType = !isEmpty(formValues) && formValues.moving_expense_type !== '';
    const isStorageExpense = this.isStorageExpense(formValues);
    const expenseNumber = expenses.length + 1;
    return (
      <div className="grid-container usa-prose">
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
        <div className="usa-grid">
          <DocumentsUploaded moveId={moveId} />
        </div>

        <div className="usa-grid expenses-container">
          <h3 className="expenses-header">Expense {expenseNumber}</h3>
          <p>
            Upload expenses one at a time.{' '}
            <Link to="/allowable-expenses" className="usa-link">
              <FontAwesomeIcon aria-hidden className="color_blue_link" icon="circle-question" />
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
                {isStorageExpense ? (
                  <div>
                    <h4 className="expenses-group-header">Dates</h4>
                    <SwaggerField
                      className="short-field expenses-form-group"
                      fieldName="storage_start_date"
                      title="Start date"
                      swagger={expenseSchema}
                      required
                    />
                    <SwaggerField
                      className="short-field expenses-form-group"
                      fieldName="storage_end_date"
                      title="End date"
                      swagger={expenseSchema}
                      required
                    />
                  </div>
                ) : (
                  <SwaggerField title="Document title" fieldName="title" swagger={moveDocSchema} required />
                )}
                <SwaggerField
                  className="short-field expense-form-element"
                  title="Amount"
                  fieldName="requested_amount_cents"
                  swagger={moveDocSchema}
                  required
                />
                <div className="expenses-uploader">
                  <p>{documentSizeLimitMsg}</p>
                  <Uploader
                    options={{ labelIdle: uploadReceipt }}
                    isPublic={isPublic}
                    onRef={(ref) => (this.uploader = ref)}
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
                {isStorageExpense && missingReceipt && (
                  <span data-testid="storage-warning">
                    <Alert type="warning">
                      If you can, go online and print a new copy of your receipt, then upload it. <br />
                      Otherwise, write and sign a statement that explains why this receipt is missing, then upload it.
                      Finance will approve or reject this expense based on your information.
                    </Alert>
                  </span>
                )}
                <div className="payment-method-radio-group-wrapper">
                  <p className="radio-group-header">How did you pay for this?</p>
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Government travel charge card (GTCC)"
                    value={paymentMethods.GTCC}
                    name="paymentMethod"
                    checked={paymentMethod === paymentMethods.GTCC}
                    onChange={this.handleRadioChange}
                  />
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Other"
                    value={paymentMethods.Other}
                    name="paymentMethod"
                    checked={paymentMethod === paymentMethods.Other}
                    onChange={this.handleRadioChange}
                  />
                  <FontAwesomeIcon
                    aria-hidden
                    className="color_blue_link"
                    icon="circle-question"
                    onClick={this.handleHowDidYouPayForThis}
                  />
                </div>
                <div className="dashed-divider" />
                <div className="radio-group-wrapper">
                  <p className="radio-group-header">Do you have more expenses to upload?</p>
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Yes"
                    value="Yes"
                    name="haveMoreExpenses"
                    checked={haveMoreExpenses === 'Yes'}
                    onChange={this.handleRadioChange}
                  />
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="No"
                    value="No"
                    name="haveMoreExpenses"
                    checked={haveMoreExpenses === 'No'}
                    onChange={this.handleRadioChange}
                  />
                </div>
              </>
            )}
            <PPMPaymentRequestActionBtns
              nextBtnLabel={nextBtnLabel}
              hasConfirmation={true}
              submitButtonsAreDisabled={this.isInvalidUploaderState() || invalid}
              submitting={submitting}
              skipHandler={this.skipHandler}
              displaySkip={expenses.length >= 1}
              saveAndAddHandler={handleSubmit(this.saveAndAddHandler)}
            />
          </form>
        </div>
      </div>
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
    expenseSchema: get(state, 'swaggerInternal.spec.definitions.CreateMovingExpenseDocumentPayload', {}),
    currentPpm: selectCurrentPPM(state) || {},
    expenses: selectPPMCloseoutDocumentsForMove(state, moveId, ['EXPENSE']),
  };
}

const mapDispatchToProps = {
  //TODO we can possibly remove selectPPMCloseoutDocumentsForMove and
  // getMoveDocumentsForMove once the document reviewer component is added
  // as it may be possible to get the number of expenses from that.
  selectPPMCloseoutDocumentsForMove,
  getMoveDocumentsForMove,
  createMovingExpenseDocument,
};

export default withContext(withLastLocation(connect(mapStateToProps, mapDispatchToProps)(ExpensesUpload)));
