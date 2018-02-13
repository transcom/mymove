import React from 'react';

function FeedbackForm({ handleChange, handleSubmit, textValue }) {
  return (
    <form onSubmit={handleSubmit}>
      <textarea
        className="feedback-form"
        onChange={handleChange}
        placeholder="Type feedback here."
        value={textValue}
      />
      <input type="submit" value="submit" />
    </form>
  );
}

export default FeedbackForm;
