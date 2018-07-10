import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm } from 'redux-form';
import DayPicker from 'react-day-picker';
import 'react-day-picker/lib/style.css';

// import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { no_op } from 'shared/utils';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';

import './DatePicker.css';

const formName = 'hhg_date_picker';
const schema = {
  properties: {
    planned_move_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Move Date',
      'x-nullable': true,
      'x-always-required': true,
    },
  },
};
const HhgDateWizardForm = reduxifyWizardForm(formName);

export class HhgDatePicker extends Component {
  state = { showInfo: false };

  constructor(props) {
    super(props);
    this.handleDayClick = this.handleDayClick.bind(this);
    this.state = {
      selectedDay: undefined,
    };
  }
  handleDayClick(day) {
    this.setState({ selectedDay: day });
    this.setState({ showInfo: true });
  }

  render() {
    const {
      pages,
      pageKey,
      error,
      currentOrders,
      serviceMemberId,
      hasSubmitSuccess,
    } = this.props;

    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders ? currentOrders : null;

    return (
      <HhgDateWizardForm
        handleSubmit={no_op}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <div className="usa-grid">
          <h4 className="date-heading">Shipment 1 (HHG)</h4>
          <h3>
            Great! Let's find a date for a moving company to move your stuff.
          </h3>
          <h4>Select a move date</h4>

          <div className="usa-width-one-third">
            <DayPicker
              onDayClick={this.handleDayClick}
              selectedDays={this.state.selectedDay}
            />
          </div>

          <div className="usa-width-two-thirds">
            {this.state.showInfo && (
              <table>
                <tbody>
                  <tr>
                    <th>Preferred Moving Dates Summary</th>
                  </tr>
                  <tr>
                    <td>Movers Packing</td>
                    <td>
                      Wed, June 6 - Thur, June 7{' '}
                      <span className="estimate">*estimated</span>
                    </td>
                  </tr>
                  <tr>
                    <td>Movers Loading Truck</td>
                    <td>Fri, June 8</td>
                  </tr>
                  <tr>
                    <td>Moving Truck in Transit</td>
                    <td>Fri, June 8 - Mon, June 11</td>
                  </tr>
                  <tr>
                    <td>Movers Delivering</td>
                    <td>
                      Tues, June 12 <span className="estimate">*estimated</span>
                    </td>
                  </tr>
                  <tr>
                    <td>Report By Date</td>
                    <td>Monday, July 16</td>
                  </tr>
                </tbody>
              </table>
            )}
          </div>
        </div>
      </HhgDateWizardForm>
    );
  }
}
HhgDatePicker.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  const props = {
    // schema: get(
    //   state,
    //   'swagger.spec.definitions.UpdateHouseholdGoodsPayload',
    //   {},
    // ),
    schema,
    formValues: getFormValues(formName)(state),
    hasSubmitSuccess: state.orders.currentOrders
      ? state.orders.hasSubmitSuccess
      : state.hhg.hasSubmitSuccess,
    ...state.hhg,
  };
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default reduxForm({ form: formName })(
  connect(mapStateToProps, mapDispatchToProps)(HhgDatePicker),
);
