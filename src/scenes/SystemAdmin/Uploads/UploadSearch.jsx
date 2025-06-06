import React, { Component } from 'react';
import { Navigate } from 'react-router-dom';
import styles from './UploadSearch.module.scss';

export class UploadSearch extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      showUpload: false,
      uploadID: '',
    };
  }

  setUploadIDinState = (e) => {
    this.setState({ uploadID: e.target.value });
  };

  redirectToShowUpload = (e) => {
    e.preventDefault();
    if (this.state.uploadID.trim()) {
      this.setState({ showUpload: true });
    } else {
      alert('Please enter a valid Upload ID.');
    }
  };

  render() {
    if (this.state.showUpload) {
      return <Navigate to={`/system/uploads/${this.state.uploadID}/show`} replace />;
    }

    return (
      <div className={styles.container}>
        <h2>Search by Upload ID</h2>
        <form onSubmit={this.redirectToShowUpload} className={styles.form}>
          <div className={styles.formGroup}>
            <label htmlFor="uploadID" className={styles.label}>
              Upload ID
            </label>
            <input
              id="uploadID"
              name="uploadID"
              type="text"
              value={this.state.uploadID}
              onChange={this.setUploadIDinState}
              className={styles.input}
              aria-required="true"
              aria-describedby="uploadIDHelp"
            />
            <small id="uploadIDHelp" className={styles.helpText}>
              Enter the unique ID associated with your upload.
            </small>
          </div>
          <button type="submit" className={styles.button} aria-label="Search for the upload">
            Search
          </button>
        </form>
      </div>
    );
  }
}

export default UploadSearch;
