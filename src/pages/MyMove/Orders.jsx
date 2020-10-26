/*  */
import { get, isEmpty } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { Field } from 'redux-form';

import {
  createOrders,
  updateOrders,
  fetchLatestOrders,
  selectActiveOrLatestOrders,
} from 'shared/Entities/modules/orders';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { withContext } from 'shared/AppContext';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import { createModifiedSchemaForOrdersTypesFlag } from 'shared/featureFlags';

import SectionWrapper from 'components/Customer/SectionWrapper';

const validateOrdersForm = validateAdditionalFields(['new_duty_station']);

const formName = 'orders_info';
const OrdersWizardForm = reduxifyWizardForm(formName, validateOrdersForm);

export class Orders extends Component {
  componentDidMount() {
    const { serviceMemberId } = this.props;
    if (!isEmpty(this.props.currentOrders)) {
      this.props.fetchLatestOrders(serviceMemberId);
    }
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
      if (isEmpty(this.props.currentOrders)) {
        return this.props.createOrders(pendingValues);
      } else {
        return this.props.updateOrders(this.props.currentOrders.id, pendingValues);
      }
    }
  };

  render() {
    const { pages, pageKey, error, currentOrders, serviceMemberId, newDutyStation, currentStation } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders ? currentOrders : null;
    const newDutyStationErrorMsg =
      newDutyStation.name === currentStation.name
        ? 'You entered the same duty station for your origin and destination. Please change one of them.'
        : '';
    const showAllOrdersTypes = this.props.context.flags.allOrdersTypes;
    const modifiedSchemaForOrdersTypesFlag = createModifiedSchemaForOrdersTypesFlag(this.props.schema);

    return (
      <OrdersWizardForm
        additionalParams={{ serviceMemberId }}
        className={formName}
        handleSubmit={this.handleSubmit}
        initialValues={initialValues}
        pageKey={pageKey}
        pageList={pages}
        readyToSubmit={!newDutyStationErrorMsg}
        serverError={error}
      >
        <h1>Tell us about your move orders</h1>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <SwaggerField
              fieldName="orders_type"
              swagger={showAllOrdersTypes ? this.props.schema : modifiedSchemaForOrdersTypesFlag}
              required
            />
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
          </div>
        </SectionWrapper>
      </OrdersWizardForm>
    );
  }
}
Orders.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
};

function mapStateToProps(state) {
  const formValues = getFormValues(formName)(state);
  const serviceMemberId = get(state, 'serviceMember.currentServiceMember.id');

  return {
    serviceMemberId: serviceMemberId,
    currentOrders: selectActiveOrLatestOrders(state),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
    formValues,
    currentStation: get(state, 'serviceMember.currentServiceMember.current_station', {}),
    newDutyStation: get(formValues, 'new_duty_station', {}),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      fetchLatestOrders,
      updateOrders,
      createOrders,
    },
    dispatch,
  );
}
export default withContext(connect(mapStateToProps, mapDispatchToProps)(Orders));
