import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import RadioButton from 'shared/RadioButton';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';

import './Expenses.css';
import { connect } from 'react-redux';
import DocumentsUploaded from './PaymentReview/DocumentsUploaded';
import withRouter from 'utils/routing';

const reviewPagePath = '/ppm-payment-review';
const nextPagePath = '/ppm-expenses';

class ExpensesLanding extends Component {
  state = {
    hasExpenses: '',
  };

  handleRadioChange = (event) => {
    this.setState({
      [event.target.name]: event.target.value,
    });
  };

  saveAndAddHandler = () => {
    const {
      moveId,
      router: { navigate },
    } = this.props;
    const { hasExpenses } = this.state;
    if (hasExpenses === 'No') {
      return navigate(`/moves/${moveId}${reviewPagePath}`);
    }
    return navigate(`/moves/${moveId}${nextPagePath}`);
  };

  render() {
    const { hasExpenses } = this.state;
    const {
      router: { navigate },
      moveId,
    } = this.props;
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
        <div className="grid-row">
          <div className="grid-col-12">
            <DocumentsUploaded moveId={moveId} />
          </div>
        </div>

        <div className="grid-row expenses-container">
          <div className="grid-col-12">
            <h3 className="expenses-header">Do you have any storage or moving expenses?</h3>
            <ul className="expenses-list">
              <li>
                <strong>Storage</strong> expenses are <strong>reimbursable</strong>.
              </li>
              <li>
                Claimable <strong>moving expenses</strong> (such as weighing fees, rental equipment, or tolls){' '}
                <strong>reduce taxes</strong> on your payment.
              </li>
            </ul>
            <Link to="/allowable-expenses" className="usa-link">
              More about expenses
            </Link>
            <div className="has-expenses-radio-group">
              <RadioButton
                inputClassName="usa-radio__input inline_radio"
                labelClassName="usa-radio__label inline_radio"
                label="Yes"
                value="Yes"
                name="hasExpenses"
                checked={hasExpenses === 'Yes'}
                onChange={this.handleRadioChange}
              />
              <RadioButton
                inputClassName="usa-radio__input inline_radio"
                labelClassName="usa-radio__label inline_radio"
                label="No"
                value="No"
                name="hasExpenses"
                checked={hasExpenses === 'No'}
                onChange={this.handleRadioChange}
              />
            </div>
            <PPMPaymentRequestActionBtns
              cancelHandler={() => {}}
              nextBtnLabel="Continue"
              saveAndAddHandler={this.saveAndAddHandler}
              finishLaterHandler={() => navigate('/')}
              submitButtonsAreDisabled={!hasExpenses}
            />
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state, { router: { params } }) {
  const moveId = params.moveId;
  return {
    moveId: moveId,
  };
}

export default withRouter(connect(mapStateToProps)(ExpensesLanding));
