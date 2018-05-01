import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { createOrders, updateOrders, loadOrders } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

const YesNoBoolean = props => {
  const {
    input: { value: rawValue, onChange },
  } = props;
  let value = rawValue;
  const yesChecked = value === true;
  const noChecked = value === false;

  const localOnChange = event => {
    if (event.target.id === 'yes') {
      onChange(true);
    } else {
      onChange(false);
    }
  };

  return (
    <Fragment>
      <input
        id="yes"
        type="radio"
        onChange={localOnChange}
        checked={yesChecked}
      />
      <label htmlFor="yes">Yes</label>
      <input
        id="no"
        type="radio"
        onChange={localOnChange}
        checked={noChecked}
      />
      <label htmlFor="no">No</label>
    </Fragment>
  );
};

const uiSchema = {
  title: 'Tell Us About Your Move Orders',
  order: [
    'orders_type',
    'issue_date',
    'report_by_date',
    'has_dependents',
    'new_duty_station',
  ],
  requiredFields: [
    'orders_type',
    'issue_date',
    'report_by_date',
    'has_dependents',
    'new_duty_station',
  ],
  custom_components: {
    has_dependents: YesNoBoolean,
    new_duty_station: DutyStationSearchBox,
  },
};
const subsetOfFields = [
  'orders_type',
  'issue_date',
  'report_by_date',
  'has_dependents',
  'new_duty_station',
];
const formName = 'orders_info';
const OrdersWizardForm = reduxifyWizardForm(formName);

export class Orders extends Component {
  componentDidMount() {
    if (this.props.currentOrders) {
      this.props.loadOrders(this.props.match.params.orderId);
    }
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    // Update if orders object already extant
    if (pendingValues) {
      const toCreateOrUpdate = pick(pendingValues, subsetOfFields);
      if (this.props.currentOrders) {
        this.props.updateOrders(this.props.currentOrders.id, toCreateOrUpdate);
      } else {
        this.props.createOrders(toCreateOrUpdate);
      }
    }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders
      ? pick(currentOrders, subsetOfFields)
      : null;
    return (
      <OrdersWizardForm
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        error={error}
        initialValues={initialValues}
        schema={this.props.schema}
        uiSchema={uiSchema}
      />
    );
  }
}
Orders.propTypes = {
  schema: PropTypes.object.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateOrders, createOrders, loadOrders },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateUpdateOrdersPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Orders);
