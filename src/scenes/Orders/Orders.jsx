import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { Field } from 'redux-form';

// import { createOrders, updateOrders } from './ducks';
import {
  createOrders,
  updateOrders,
  showServiceMemberOrders,
  getLatestOrdersLabel,
  // selectOrders,
} from 'shared/Entities/modules/orders';
import { getRequestStatus } from 'shared/Swagger/selectors';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { withContext } from 'shared/AppContext';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';

const validateOrdersForm = validateAdditionalFields(['new_duty_station']);

const formName = 'orders_info';
const OrdersWizardForm = reduxifyWizardForm(formName, validateOrdersForm);

export class Orders extends Component {
  componentDidMount() {
    const { serviceMemberId } = this.props;
    this.props.showServiceMemberOrders(serviceMemberId);
  }

  handleSubmit = () => {
    const pendingValues = Object.assign({}, this.props.formValues);

    // Update if orders object already extant
    if (pendingValues) {
      pendingValues['service_member_id'] = this.props.serviceMemberId;
      pendingValues['new_duty_station_id'] = pendingValues.new_duty_station.id;
      pendingValues['has_dependents'] = pendingValues.has_dependents || false;
      pendingValues['spouse_has_pro_gear'] =
        (pendingValues.has_dependents && pendingValues.spouse_has_pro_gear) || false;
      if (this.props.currentOrders) {
        return this.props.update(this.props.currentOrders.id, pendingValues);
      } else {
        return this.props.create(pendingValues);
      }
    }
  };

  render() {
    const {
      pages,
      pageKey,
      error,
      // currentOrders,
      orders,
      serviceMemberId,
      newDutyStation,
      currentStation,
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = orders ? orders : null;
    // const initialValues = currentOrders ? currentOrders : null;
    const newDutyStationErrorMsg =
      newDutyStation.name === currentStation.name
        ? 'You entered the same duty station for your origin and destination. Please change one of them.'
        : '';
    return (
      <OrdersWizardForm
        additionalParams={{ serviceMemberId }}
        className={formName}
        handleSubmit={this.handleSubmit}
        initialValues={initialValues}
        readyToSubmit={!newDutyStationErrorMsg}
        pageKey={pageKey}
        pageList={pages}
        serverError={error}
      >
        <h1 className="sm-heading">Tell us about your move orders</h1>
        <SwaggerField fieldName="orders_type" swagger={this.props.schema} required />
        <SwaggerField fieldName="issue_date" swagger={this.props.schema} required />
        <div style={{ marginTop: '0.25rem' }}>
          <span className="usa-hint">Date your orders were issued.</span>
        </div>
        <SwaggerField fieldName="report_by_date" swagger={this.props.schema} required />
        <SwaggerField fieldName="has_dependents" swagger={this.props.schema} component={YesNoBoolean} />
        <Field
          name="new_duty_station"
          component={DutyStationSearchBox}
          errorMsg={newDutyStationErrorMsg}
          title="New duty station"
        />
      </OrdersWizardForm>
    );
  }
}
Orders.propTypes = {
  schema: PropTypes.object.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: PropTypes.object,
  error: PropTypes.object,
};

function mapStateToProps(state) {
  const formValues = getFormValues(formName)(state);
  const showOrdersRequest = getRequestStatus(state, getLatestOrdersLabel);
  // const ordersId = get(state, `entities.orders`);

  return {
    serviceMemberId: get(state, 'serviceMember.currentServiceMember.id'),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
    formValues,
    currentOrders: state.orders.currentOrders,
    // orders: get(state, 'entities.orders'),
    // orders: selectOrders(state, ordersId),
    currentStation: get(state, 'serviceMember.currentServiceMember.current_station', {}),
    newDutyStation: get(formValues, 'new_duty_station', {}),
    loadDependenciesHasSuccess: showOrdersRequest.isSuccess,
    loadDependenciesHasError: showOrdersRequest.error,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      showServiceMemberOrders,
      update: updateOrders,
      create: createOrders,
    },
    dispatch,
  );
}
export default withContext(connect(mapStateToProps, mapDispatchToProps)(Orders));
