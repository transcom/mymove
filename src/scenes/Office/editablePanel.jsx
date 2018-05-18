import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Alert from 'shared/Alert';
import { EditablePanel } from 'shared/EditablePanel';

export default function editablePanel(DisplayComponent, EditComponent) {
  const Wrapper = class extends Component {
    constructor(props) {
      super(props);
      this.state = {
        isEditable: false,
      };
      // TODO: Figure out why bind is still needed when ostensibly it's not
      this.save = this.save.bind(this);
    }

    save = () => {
      const args = this.props.getUpdateArgs();
      this.props.update(...args);
      this.toggleEditable();
    };

    toggleEditable = () => {
      this.setState({
        isEditable: !this.state.isEditable,
      });
    };

    render() {
      const isEditable = this.state.isEditable || this.props.isUpdating;
      const Content = isEditable ? EditComponent : DisplayComponent;

      return (
        <React.Fragment>
          {this.props.hasError && (
            <Alert type="error" heading="An error occurred">
              There was an error: <em>{this.props.errorMessage}</em>.
            </Alert>
          )}
          <EditablePanel
            title={this.props.title}
            onSave={this.save}
            onToggle={this.toggleEditable}
            isEditable={isEditable}
          >
            <Content
              values={this.props.displayValues}
              schema={this.props.schema}
            />
          </EditablePanel>
        </React.Fragment>
      );
    }
  };

  Wrapper.propTypes = {
    schema: PropTypes.object.isRequired,
    displayValues: PropTypes.object.isRequired,
    update: PropTypes.func.isRequired,
    moveId: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
    isUpdating: PropTypes.bool.isRequired,
  };

  return Wrapper;
}
