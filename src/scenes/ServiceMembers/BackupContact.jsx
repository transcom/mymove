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
  deleteBackupContact,
} from './ducks';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

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

function hasEmptyValues(values) {
  if (!values) {
    return true;
  }

  const emptyableFields = ['name', 'email', 'telephone'];
  let allZero = true;
  emptyableFields.forEach(fieldName => {
    const value = values[fieldName]; // eslint-disable-line security/detect-object-injection
    if (value !== undefined && value !== '' && value !== null) {
      allZero = false;
    }
  });

  return allZero;
}

const BackupContactWizardForm = reduxifyWizardForm(formName);

export class BackupContact extends Component {
  handleSubmit = () => {
    const pendingValues = this.props.values;

    if (pendingValues.telephone === '') {
      pendingValues.telephone = null;
    }

    if (pendingValues.permission === undefined) {
      pendingValues.permission = NonePermission;
    }

    if (pendingValues && !hasEmptyValues(pendingValues)) {
      if (this.props.currentBackupContacts.length > 0) {
        // update existing
        const oldOne = this.props.currentBackupContacts[0];
        return this.props.updateBackupContact(oldOne.id, pendingValues);
      } else {
        return this.props.createBackupContact(
          this.props.match.params.serviceMemberId,
          pendingValues,
        );
      }
    }

    // If we have empty values and an existing BC, delete it.
    if (
      hasEmptyValues(pendingValues) &&
      this.props.currentBackupContacts.length > 0
    ) {
      const oldOne = this.props.currentBackupContacts[0];
      return this.props.deleteBackupContact(oldOne.id);
    }
  };

  render() {
    const { pages, pageKey, error, schema } = this.props;

    const isSkippable = hasEmptyValues(this.props.values);

    // eslint-disable-next-line
    var [contact1, contact2] = this.props.currentBackupContacts; // contact2 will be used when we implement saving two backup contacts.

    // initialValues has to be null until there are values from the action since only the first values are taken
    const firstInitialValues = contact1
      ? pick(contact1, ['name', 'email', 'telephone', 'permission'])
      : null;

    return (
      <BackupContactWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={firstInitialValues}
        skippable={isSkippable}
      >
        <h1 className="sm-heading">
          Backup Contact (Trusted Agents){' '}
          <span className="optional">optional</span>
        </h1>
        <p>
          If we can't reach you, who can we contact (such as spouse or parent)*?
        </p>

        <SwaggerField fieldName="name" swagger={schema} required />
        <SwaggerField fieldName="email" swagger={schema} required />
        <SwaggerField fieldName="telephone" swagger={schema} />

        {/* TODO: Uncomment line below after backup contact auth is implemented.  */}
        {/* <Field name="permission" component={permissionsField} /> */}
        <p>
          * Any person you assign as a backup or trusted agent must be 18 years
          of age or older.
        </p>
      </BackupContactWizardForm>
    );
  }
}
BackupContact.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      updateServiceMember,
      createBackupContact,
      updateBackupContact,
      deleteBackupContact,
    },
    dispatch,
  );
}
function mapStateToProps(state) {
  return {
    currentBackupContacts: state.serviceMember.currentBackupContacts,
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
