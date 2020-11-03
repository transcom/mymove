/* eslint-disable react/forbid-prop-types */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, Field } from 'redux-form';

import {
  createOrders as createOrdersAction,
  updateOrders as updateOrdersAction,
  fetchLatestOrders as fetchLatestOrdersAction,
  selectActiveOrLatestOrders,
} from 'shared/Entities/modules/orders';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { withContext } from 'shared/AppContext';
import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';
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
    const { serviceMemberId, currentOrders, fetchLatestOrders } = this.props;
    if (!isEmpty(currentOrders)) {
      fetchLatestOrders(serviceMemberId);
    }
  }

  handleSubmit = () => {
    const { formValues, serviceMemberId, currentOrders, createOrders, updateOrders } = this.props;
    const pendingValues = { ...formValues };

    // Update if orders object already extant
    if (pendingValues) {
      pendingValues.service_member_id = serviceMemberId;
      pendingValues.new_duty_station_id = pendingValues.new_duty_station.id;
      pendingValues.has_dependents = pendingValues.has_dependents || false;
      pendingValues.spouse_has_pro_gear = (pendingValues.has_dependents && pendingValues.spouse_has_pro_gear) || false;

      if (isEmpty(currentOrders)) {
        return createOrders(pendingValues);
      }

      return updateOrders(currentOrders.id, pendingValues);
    }

    return null;
  };

  render() {
    const {
      schema,
      context,
      pages,
      pageKey,
      error,
      currentOrders,
      serviceMemberId,
      newDutyStation,
      currentStation,
    } = this.props;

    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders || null;
    const newDutyStationErrorMsg =
      newDutyStation.name === currentStation.name
        ? 'You entered the same duty station for your origin and destination. Please change one of them.'
        : '';
    const showAllOrdersTypes = context.flags.allOrdersTypes;
    const modifiedSchemaForOrdersTypesFlag = createModifiedSchemaForOrdersTypesFlag(schema);

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
              swagger={showAllOrdersTypes ? schema : modifiedSchemaForOrdersTypesFlag}
              required
            />
            <SwaggerField fieldName="issue_date" swagger={schema} required />
            <div style={{ marginTop: '0.25rem' }}>
              <span className="usa-hint">Date your orders were issued.</span>
            </div>
            <SwaggerField fieldName="report_by_date" swagger={schema} required />
            <SwaggerField fieldName="has_dependents" swagger={schema} component={YesNoBoolean} />
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
  serviceMemberId: PropTypes.string.isRequired,
  currentOrders: PropTypes.object,
  fetchLatestOrders: PropTypes.func,
  createOrders: PropTypes.func,
  updateOrders: PropTypes.func,
  formValues: PropTypes.object,
  pages: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  newDutyStation: PropTypes.object,
  currentStation: PropTypes.object,
};

Orders.defaultProps = {
  error: null,
  currentOrders: null,
  fetchLatestOrders: () => {},
  createOrders: () => {},
  updateOrders: () => {},
  formValues: {},
  newDutyStation: {},
  currentStation: {},
};

function mapStateToProps(state) {
  const formValues = getFormValues(formName)(state);
  const serviceMemberId = get(state, 'serviceMember.currentServiceMember.id');

  return {
    serviceMemberId,
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
      fetchLatestOrders: fetchLatestOrdersAction,
      updateOrders: updateOrdersAction,
      createOrders: createOrdersAction,
    },
    dispatch,
  );
}
export default withContext(connect(mapStateToProps, mapDispatchToProps)(Orders));
