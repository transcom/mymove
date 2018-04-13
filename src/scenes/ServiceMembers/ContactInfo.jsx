import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { updateServiceMember, loadServiceMember } from './ducks';
import WizardPage from 'shared/WizardPage';
import { no_op } from 'shared/utils';
import { pick } from 'lodash';

import { reduxifyForm } from 'shared/JsonSchemaForm';
import './ContactInfo.css';

// import { Field, reduxForm } from 'redux-form';
// import validator from 'shared/JsonSchemaForm/validator';
// function Phone({ name, isRequired }) {
//   const validators = [
//     // validator.patternMatches(
//     //   '^[2-9]d{2}-d{3}-d{4}$',
//     //   'Number must have 10 digits.',
//     // ),
//   ];
//   if (isRequired) validators.push(validator.isRequired);
//   return (
//     <Field
//       name={name}
//       component="input"
//       type="text"
//       normalize={validator.normalizePhone}
//       validate={validators}
//     />
//   );
// }
// function Form(props) {
//   const { handleSubmit } = props;
//   return (
//     <div onClick={handleSubmit}>
//       <form>
//         <h2>Your Contact Info</h2>
//         <label>
//           Best contact phone
//           <Phone name="telephone" isRequired={true} />
//         </label>
//         <label>
//           Alternate phone <span className="label-optional">Optional</span>
//           <Phone name="secondary_telephone" />
//         </label>
//         <label>
//           Email
//           <Field
//             name="personal_email"
//             component="input"
//             type="text"
//             validate={validator.isRequired}
//           />
//         </label>
//       </form>
//     </div>
//   );
// }

// const CurrentForm = reduxForm({
//   form: 'service_member_contact_info',
// })(Form);
const CurrentForm = reduxifyForm('service_member_contact_info');
export class ContactInfo extends Component {
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    this.props.updateServiceMember({});
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
      userEmail,
    } = this.props;
    const isValid = this.refs.currentForm && !this.refs.currentForm.valid;
    const isDirty = this.refs.currentForm && this.refs.currentForm.dirty;
    const subsetOfFields = [
      'telephone',
      'secondary_telephone',
      'personal_email',
      'phone_is_preferred',
      'secondary_phone_is_preferred',
      'email_is_preferred',
    ];
    const initialValues = currentServiceMember
      ? pick(currentServiceMember, subsetOfFields)
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
          className="contact-info"
          handleSubmit={no_op}
          schema={this.props.schema}
          uiSchema={this.props.uiSchema}
          showSubmit={false}
          initialValues={initialValues}
        />
      </WizardPage>
    );
  }
}
ContactInfo.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    userEmail: state.user.email,
    schema: {},
    ...state.serviceMember,
    uiSchema: {
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
            'secondary_phone_is_preferred',
            'email_is_preferred',
          ],
        },
      },
    },
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateServiceMemberPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(ContactInfo);
