import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { createOrders, updateOrders, loadOrders } from './ducks';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

// import './Orders.css';

const uiSchema = {
  title: 'Your Move Orders',
  order: ['orders_type', 'issue_date', 'report_by_date', 'has_dependents'],
  requiredFields: [
    'orders_type',
    'issue_date',
    'report_by_date',
    'has_dependents',
  ],
  // groups: {
  //   has_dependents: {
  //     title: 'Are dependents included in your orders?',
  //     fields: [
  //       'has_dependents',
  //     ],
  //   },
  // },
};
const subsetOfFields = [
  'orders_type',
  'issue_date',
  'report_by_date',
  'has_dependents',
];
const formName = 'orders_info';
const CurrentForm = reduxifyForm(formName);

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
    // let prefSelected = false;
    // if (this.refs.currentForm && this.refs.currentForm.values) {
    //   prefSelected = Boolean(
    //     this.refs.currentForm.values.phone_is_preferred ||
    //       this.refs.currentForm.values.text_message_is_preferred ||
    //       this.refs.currentForm.values.email_is_preferred,
    //   );
    // }
    const isValid = this.refs.currentForm && this.refs.currentForm.valid;
    const isDirty = this.refs.currentForm && this.refs.currentForm.dirty;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders
      ? pick(currentOrders, subsetOfFields)
      : null;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={isDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <CurrentForm
          ref="currentForm"
          className={formName}
          handleSubmit={no_op}
          schema={this.props.schema}
          uiSchema={uiSchema}
          showSubmit={false}
          initialValues={initialValues}
        />
      </WizardPage>
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
  return bindActionCreators({ updateOrders, loadOrders }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.OrdersPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Orders);
