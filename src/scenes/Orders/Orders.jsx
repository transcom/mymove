import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateOrders, loadOrders } from './ducks';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import './Orderss.css';

const uiSchema = {
  title: 'Your Contact Info',
  order: [
    'telephone',
    'secondary_telephone',
    'personal_email',
    'contact_preferences',
  ],
  requiredFields: ['telephone', 'personal_email'],
  groups: {
    contact_preferences: {
      title: 'Preferred contact method during your move:',
      fields: [
        'phone_is_preferred',
        'text_message_is_preferred',
        'email_is_preferred',
      ],
    },
  },
};
const subsetOfFields = [
  'telephone',
  'secondary_telephone',
  'personal_email',
  'phone_is_preferred',
  'text_message_is_preferred',
  'email_is_preferred',
];
const formName = 'orders_info';
const CurrentForm = reduxifyForm(formName);

export class Orders extends Component {
  componentDidMount() {
    this.props.loadOrders(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      this.props.updateOrders(patch);
    }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
      userEmail,
    } = this.props;
    let prefSelected = false;
    if (this.refs.currentForm && this.refs.currentForm.values) {
      prefSelected = Boolean(
        this.refs.currentForm.values.phone_is_preferred ||
          this.refs.currentForm.values.text_message_is_preferred ||
          this.refs.currentForm.values.email_is_preferred,
      );
    }
    const isValid =
      this.refs.currentForm && this.refs.currentForm.valid && prefSelected;
    const isDirty = this.refs.currentForm && this.refs.currentForm.dirty;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders
      ? pick(currentOrders, subsetOfFields)
      : null;
    if (initialValues && !initialValues.personal_email)
      initialValues.personal_email = userEmail;
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
  userEmail: PropTypes.string.isRequired,
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
    userEmail: state.user.email,
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateOrdersPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Orders);
