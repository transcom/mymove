import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import {
  updateServiceMember,
  createBackupContact,
  updateBackupContact,
} from './ducks';
import {
  renderField,
  recursivelyAnnotateRequiredFields,
} from 'shared/JsonSchemaForm';
import { reduxForm } from 'redux-form';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import './BackupContact.css';

const NonePermission = 'NONE';
// const ViewPermission = 'VIEW';
// const EditPermission = 'EDIT';

// TODO: Uncomment field below after backup contact auth is implemented.
// const permissionsField = props => {
//   const {
//     input: { value: rawValue, onChange },
//   } = props;
//   let value;
//   if (![NonePermission, ViewPermission, EditPermission].includes(rawValue)) {
//     value = NonePermission;
//   } else {
//     value = rawValue;
//   }

//   const localOnChange = event => {
//     if (event.target.id === 'authorizeAgent') {
//       if (event.target.checked && value === NonePermission) {
//         onChange(ViewPermission);
//       } else if (!event.target.checked) {
//         onChange(NonePermission);
//       }
//     } else if (event.target.id === 'aaChoiceView') {
//       onChange(ViewPermission);
//     } else if (event.target.id === 'aaChoiceEdit') {
//       onChange(EditPermission);
//     }
//   };

//   const authorizedChecked = value !== NonePermission;
//   const viewChecked = value === ViewPermission;
//   const editChecked = value === EditPermission;

//   return (
//     <Fragment>
//       <input
//         id="authorizeAgent"
//         type="checkbox"
//         onChange={localOnChange}
//         checked={authorizedChecked}
//       />
//       <label htmlFor="authorizeAgent">I authorize this person to:</label>
//       <input
//         id="aaChoiceView"
//         type="radio"
//         onChange={localOnChange}
//         checked={viewChecked}
//         disabled={!authorizedChecked}
//       />
//       <label
//         htmlFor="aaChoiceView"
//         className={authorizedChecked ? '' : 'disabled'}
//       >
//         Sign for pickup or delivery in my absence, and view move details in this
//         app.
//       </label>
//       <input
//         id="aaChoiceEdit"
//         type="radio"
//         onChange={localOnChange}
//         checked={editChecked}
//         disabled={!authorizedChecked}
//       />
//       <label
//         htmlFor="aaChoiceEdit"
//         className={authorizedChecked ? '' : 'disabled'}
//       >
//         Represent me in all aspects of this move (this person will be invited to
//         login and will be authorized with power of attorney on your behalf).
//       </label>
//     </Fragment>
//   );
// };

const formName = 'service_member_backup_contact';

class ContactForm extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  static getDerivedStateFromProps(nextProps, prevState) {
    const { valid, dirty, updateValidDirty } = nextProps;
    updateValidDirty(valid, dirty);

    return prevState;
  }

  render() {
    const { schema } = this.props;
    recursivelyAnnotateRequiredFields(schema);
    const fields = schema.properties || {};

    return (
      <form>
        <h1 className="sm-heading">Backup Contact</h1>
        <p>
          If we can't reach you, who can we contact (such as spouse or parent)?
        </p>

        {renderField('name', fields, '')}
        {renderField('email', fields, '')}
        {renderField('telephone', fields, '')}

        {/* TODO: Uncomment line below after backup contact auth is implemented.  */}
        {/* <Field name="permission" component={permissionsField} /> */}
      </form>
    );
  }
}

const validateContact = (values, form) => {
  let requiredErrors = {};
  /* eslint-disable security/detect-object-injection */
  ['name', 'email'].forEach(requiredFieldName => {
    if (
      values[requiredFieldName] === undefined ||
      values[requiredFieldName] === ''
    ) {
      requiredErrors[requiredFieldName] = 'Required.';
    }
  });
  /* eslint-enable security/detect-object-injection */
  return requiredErrors;
};

ContactForm = reduxForm({ form: formName, validate: validateContact })(
  ContactForm,
);

export class BackupContact extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isValid: true,
      isDirty: false,
    };
  }

  handleSubmit = () => {
    const pendingValues = this.props.values;

    if (pendingValues.telephone === '') {
      pendingValues.telephone = null;
    }

    if (pendingValues.permission === undefined) {
      pendingValues.permission = NonePermission;
    }

    if (pendingValues) {
      if (this.props.currentBackupContacts.length > 0) {
        // update existing
        const oldOne = this.props.currentBackupContacts[0];
        this.props.updateBackupContact(oldOne.id, pendingValues);
      } else {
        this.props.createBackupContact(
          this.props.match.params.serviceMemberId,
          pendingValues,
        );
      }
    }
  };

  updateValidDirty = (isValid, isDirty) => {
    this.setState({
      isValid,
      isDirty,
    });
  };

  render() {
    const { pages, pageKey, hasSubmitSuccess, error } = this.props;
    const isValid = this.state.isValid;
    const isDirty = this.state.isDirty;

    // eslint-disable-next-line
    var [contact1, contact2] = this.props.currentBackupContacts; // contact2 will be used when we implement saving two backup contacts.

    // initialValues has to be null until there are values from the action since only the first values are taken
    const firstInitialValues = contact1
      ? pick(contact1, ['name', 'email', 'telephone', 'permission'])
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
        <ContactForm
          ref="currentForm"
          updateValidDirty={this.updateValidDirty}
          initialValues={firstInitialValues}
          handleSubmit={no_op}
          schema={this.props.schema}
        />
      </WizardPage>
    );
  }
}
BackupContact.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      updateServiceMember,
      createBackupContact,
      updateBackupContact,
    },
    dispatch,
  );
}
function mapStateToProps(state) {
  return {
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    hasSubmitSuccess:
      state.serviceMember.createBackupContactSuccess ||
      state.serviceMember.updateBackupContactSuccess,
    error: state.serviceMember.error,
    schema: get(
      state,
      'swagger.spec.definitions.CreateServiceMemberBackupContactPayload',
      {},
    ),
    values: getFormValues(formName)(state),
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(BackupContact);
