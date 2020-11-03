/* eslint-disable react/no-unused-prop-types */
/* eslint-disable react/forbid-prop-types */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { isEmpty } from 'lodash';
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
import { WizardPage } from 'shared/WizardPage/index';
import { MatchShape, HistoryShape, PageKeyShape, PageListShape } from 'types/customerShapes';

export class Orders extends Component {
  componentDidMount() {
    // TODO
    const { serviceMemberId, currentOrders, fetchLatestOrders } = this.props;
    if (!isEmpty(currentOrders)) {
      fetchLatestOrders(serviceMemberId);
    }
  }

  render() {
    const {
      // context,
      // error,
      currentStation,
      match,
      pages,
      pageKey,
      history,
      serviceMemberId,
      currentOrders,
      createOrders,
      updateOrders,
    } = this.props;

    const submitOrders = (values) => {
      const pendingValues = {
        ...values,
        service_member_id: serviceMemberId,
        new_duty_station_id: values.new_duty_station.id,
        has_dependents: values.has_dependents === 'yes',
        spouse_has_pro_gear: false, // TODO - this input seems to be deprecated?
      };

      if (currentOrders?.id) {
        return updateOrders(currentOrders.id, pendingValues);
      }

      return createOrders(pendingValues);
    };

    const initialValues = {
      orders_type: currentOrders?.orders_type || '',
      issue_date: currentOrders?.issue_date || '',
      report_by_date: currentOrders?.report_by_date || '',
      has_dependents: currentOrders?.has_dependents ? 'yes' : 'no', // TODO - radio is not prefilling
      new_duty_station: currentOrders?.new_duty_station || null,
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
      <Formik initialValues={initialValues} validateOnMount validationSchema={ordersInfoSchema} onSubmit={submitOrders}>
        {({ isValid, dirty, handleSubmit }) => (
          <WizardPage
            canMoveNext={dirty && isValid}
            match={match}
            pageKey={pageKey}
            pageList={pages}
            push={history.push}
            handleSubmit={handleSubmit}
            dirty={dirty}
          >
            <h1>Tell us about your move orders</h1>
            <SectionWrapper>
              <OrdersInfoForm />
            </SectionWrapper>
          </WizardPage>
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
  match: MatchShape.isRequired,
  history: HistoryShape.isRequired,
  pages: PageListShape,
  pageKey: PageKeyShape,
};

Orders.defaultProps = {
  currentOrders: null,
  fetchLatestOrders: () => {},
  createOrders: () => {},
  updateOrders: () => {},
  currentStation: {},
  pages: [],
  pageKey: '',
};

const mapStateToProps = (state) => ({
  serviceMemberId: state.serviceMember?.currentServiceMember?.id,
  currentOrders: selectActiveOrLatestOrders(state),
  currentStation: state.serviceMember?.currentServiceMember?.current_station || {},
});

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
  updateOrders: updateOrdersAction,
  createOrders: createOrdersAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(Orders));
