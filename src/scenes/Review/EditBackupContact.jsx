import React, { Component } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import scrollToTop from 'shared/scrollToTop';

import { push } from 'connected-react-router';
import { reduxForm } from 'redux-form';

import { updateBackupContact as updateBackupContactAction } from 'store/entities/actions';
import { patchBackupContact, getResponseError } from 'services/internalApi';
import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import SaveCancelButtons from './SaveCancelButtons';
import './Review.css';
import profileImage from './images/profile.png';
import { selectBackupContacts } from 'store/entities/selectors';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

const editBackupContactFormName = 'edit_backup_contact';

let EditBackupContactForm = (props) => {
  const { schema, handleSubmit, submitting, valid } = props;
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <img src={profileImage} alt="" />
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 0,
                marginTop: 20,
              }}
            >
              Backup Contact
            </h1>
            <hr />
            <h3>Edit Backup Contact:</h3>
            <p>Any person you assign as a backup contact must be 18 years of age or older.</p>
            <SwaggerField fieldName="name" swagger={schema} required />
            <SwaggerField fieldName="email" swagger={schema} required />
            <SwaggerField fieldName="telephone" swagger={schema} />
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};
EditBackupContactForm = reduxForm({
  form: editBackupContactFormName,
})(EditBackupContactForm);

class EditBackupContact extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  componentDidMount() {
    scrollToTop();
  }

  updateContact = (fieldValues) => {
    const { updateBackupContact, setFlashMessage } = this.props;

    if (fieldValues.telephone === '') {
      fieldValues.telephone = null;
    }

    return patchBackupContact(fieldValues)
      .then((response) => {
        // Update in Redux
        updateBackupContact(response);

        setFlashMessage('EDIT_BACKUP_CONTACT_SUCCESS', 'success', '', 'Your changes have been saved.');
        this.props.history.goBack();
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update backup contact due to server error');
        this.setState({
          errorMessage,
        });

        scrollToTop();
      });
  };

  render() {
    const { schema, backupContacts } = this.props;
    const { errorMessage } = this.state;

    let backupContact = null;
    if (backupContacts) {
      backupContact = backupContacts[0];
    }

    return (
      <div className="usa-grid">
        {errorMessage && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {errorMessage}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditBackupContactForm initialValues={backupContact} onSubmit={this.updateContact} schema={schema} />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    backupContacts: selectBackupContacts(state),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberBackupContactPayload', {}),
  };
}

const mapDispatchToProps = {
  push,
  updateBackupContact: updateBackupContactAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(EditBackupContact);
