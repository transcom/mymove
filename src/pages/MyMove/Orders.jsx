/* eslint-disable react/forbid-prop-types */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { Formik } from 'formik';
import * as Yup from 'yup';

import {
  createOrders as createOrdersAction,
  updateOrders as updateOrdersAction,
  fetchLatestOrders as fetchLatestOrdersAction,
  selectActiveOrLatestOrders,
} from 'shared/Entities/modules/orders';
import { withContext } from 'shared/AppContext';
import SectionWrapper from 'components/Customer/SectionWrapper';
import OrdersInfoForm from 'components/Customer/OrdersInfoForm/OrdersInfoForm';

export class Orders extends Component {
  componentDidMount() {
    // TODO
    const { serviceMemberId, currentOrders, fetchLatestOrders } = this.props;
    if (!isEmpty(currentOrders)) {
      fetchLatestOrders(serviceMemberId);
    }
  }

  handleSubmit = (values) => {
    const { serviceMemberId, currentOrders, createOrders, updateOrders } = this.props;
    const pendingValues = { ...values };

    // TODO
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
      // context,
      // pages,
      // pageKey,
      // error,
      // currentOrders,
      // serviceMemberId,
      // newDutyStation,
      currentStation,
    } = this.props;

    // initialValues has to be null until there are values from the action since only the first values are taken
    // TODO - initialize values based on currentOrders
    // const initialValues = currentOrders || null;
    const initialValues = {
      orders_type: '', // required
      issue_date: '', // required
      report_by_date: '', // required
      has_dependents: '', // required
      new_duty_station: null,
    };

    // TODO - orders types feature flag
    // const showAllOrdersTypes = context.flags.allOrdersTypes;
    const ordersTypeOptions = [
      { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
      { key: 'RETIREMENT', value: 'Retirement' },
      { key: 'SEPARATION', value: 'Separation' },
    ];

    const ordersInfoSchema = Yup.object().shape({
      orders_type: Yup.mixed()
        .oneOf(ordersTypeOptions.map((i) => i.key))
        .required('Required'),
      issue_date: Yup.date().required('Required'),
      report_by_date: Yup.date().required('Required'),
      has_dependents: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
      new_duty_station: Yup.object()
        .shape({
          name: Yup.string().notOneOf(
            [currentStation?.name],
            'You entered the same duty station for your origin and destination. Please change one of them.',
          ),
        })
        .nullable()
        .required('Required'),
    });

    return (
      <Formik initialValues={initialValues} validateOnMount validationSchema={ordersInfoSchema}>
        {() => (
          <>
            <h1>Tell us about your move orders</h1>
            <SectionWrapper>
              <OrdersInfoForm currentStation={currentStation} />
            </SectionWrapper>
          </>
        )}
      </Formik>
    );
  }
}

Orders.propTypes = {
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
  currentStation: PropTypes.object,
};

Orders.defaultProps = {
  currentOrders: null,
  fetchLatestOrders: () => {},
  createOrders: () => {},
  updateOrders: () => {},
  currentStation: {},
};

function mapStateToProps(state) {
  const serviceMemberId = get(state, 'serviceMember.currentServiceMember.id');

  return {
    serviceMemberId,
    currentOrders: selectActiveOrLatestOrders(state),
    currentStation: get(state, 'serviceMember.currentServiceMember.current_station', {}),
  };
}

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
  updateOrders: updateOrdersAction,
  createOrders: createOrdersAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(Orders));
