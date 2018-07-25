import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm } from 'redux-form';
import DayPicker from 'react-day-picker';
import 'react-day-picker/lib/style.css';

import {
  createOrUpdateShipment,
  selectShipment,
} from 'shared/Entities/modules/shipments';
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
const HHGDateWizardForm = reduxifyWizardForm(formName);

export class HHGDatePicker extends Component {
  handleSubmit = () => {
    debugger;
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    createOrUpdateShipment(moveId, shipment);
  };

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
      serviceMemberId,
      hasSubmitSuccess,
    } = this.props;

    // initialValues has to be null until there are values from the action since only the first values are taken

    return (
      <HHGDateWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={false}
        serverError={error}
        // initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <div className="usa-grid">
          <h3>Shipment 1 (HHG)</h3>
          <h2 className="instruction-heading">
            Great! Let's find a date for a moving company to move your stuff.
          </h2>
          <h3>Select a move date</h3>

          <div className="usa-width-one-third">
            <DayPicker
              onDayClick={this.handleDayClick}
              selectedDays={this.state.selectedDay}
            />
          </div>

          <div className="usa-width-two-thirds">
            {this.state.showInfo && (
              <table className="Todo-phase2">
                <tbody>
                  <tr>
                    <th className="Todo-phase2">
                      Preferred Moving Dates Summary
                    </th>
                  </tr>
                  <tr>
                    <td>Movers Packing</td>
                    <td className="Todo-phase2">
                      Wed, June 6 - Thur, June 7{' '}
                      <span className="estimate">*estimated</span>
                    </td>
                  </tr>
                  <tr>
                    <td>Movers Loading Truck</td>
                    <td className="Todo-phase2">Fri, June 8</td>
                  </tr>
                  <tr>
                    <td>Moving Truck in Transit</td>
                    <td className="Todo-phase2">Fri, June 8 - Mon, June 11</td>
                  </tr>
                  <tr>
                    <td>Movers Delivering</td>
                    <td className="Todo-phase2">
                      Tues, June 12 <span className="estimate">*estimated</span>
                    </td>
                  </tr>
                  <tr>
                    <td>Report By Date</td>
                    <td className="Todo-phase2">Monday, July 16</td>
                  </tr>
                </tbody>
              </table>
            )}
          </div>
        </div>
      </HHGDateWizardForm>
    );
  }
}
HHGDatePicker.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
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
    move: get(state, 'moves.currentMove', {}),
    shipment: {},
  };
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdateShipment }, dispatch);
}

export default reduxForm({ form: formName })(
  connect(mapStateToProps, mapDispatchToProps)(HHGDatePicker),
);
