import React, { Component } from 'react';

const PanelContext = React.createContext();

const PanelConsumer = props => (
  <PanelContext.Consumer {...props}>
    {context => {
      if (!context) {
        throw new Error(
          `Panel compound components cannot be rendered outside the Panel component`,
        );
      }
      return props.children(context);
    }}
  </PanelContext.Consumer>
);

class Panel extends Component {
  handleEditClick = () => {
    this.setState(() => ({ isEditing: true }));
  };

  handleCancel = e => {
    e.preventDefault();
    this.setState(() => ({ isEditing: false }));
  };

  handleSave = e => {
    e.preventDefault();
  };

  state = {
    isEditing: false,
    onEditClick: this.handleEditClick,
    onCancel: this.handleCancel,
    onSave: this.handleSave,
  };

  static HalfRow = ({ children }) => (
    <div className="usa-width-one-half">{children}</div>
  );

  static Title = ({ children, editLabel, editEnabled = true }) => (
    <PanelConsumer>
      {({ isEditing, onEditClick, onCancel, onSave }) => (
        <div className="editable-panel-header $">
          <div className="title">{children}</div>
          {editEnabled &&
            !isEditing && (
              <a className="editable-panel-edit" onClick={onEditClick}>
                {editLabel}
              </a>
            )}
        </div>
      )}
    </PanelConsumer>
  );

  static Header = ({ children }) => (
    <div className="column-head">{children}</div>
  );

  static Subheader = ({ children }) => (
    <div className="column-subhead">{children}</div>
  );

  static Content = ({ children }) => (
    <div className="editable-panel-content">{children}</div>
  );

  static CancelButton = ({ children }) => (
    <PanelConsumer>
      {({ onCancel }) => (
        <button
          className="usa-button-secondary editable-panel-cancel"
          onClick={onCancel}
        >
          {children}
        </button>
      )}
    </PanelConsumer>
  );

  static SaveButton = ({ children }) => (
    <PanelConsumer>
      {({ onSave }) => (
        <button className="usa-button editable-panel-save" onClick={onSave}>
          {children}
        </button>
      )}
    </PanelConsumer>
  );

  get getStateAndMethods() {
    return {
      ...this.state,
      onEditClick: this.handleEditClick,
      onCancel: this.handleCancel,
      onSave: this.handleSave,
    };
  }

  render() {
    const { className } = this.props;
    return (
      <div className={`editable-panel ${className}`}>
        <PanelContext.Provider value={this.state}>
          {this.props.children(this.getStateAndMethods)}
        </PanelContext.Provider>
      </div>
    );
  }
}

export default Panel;
