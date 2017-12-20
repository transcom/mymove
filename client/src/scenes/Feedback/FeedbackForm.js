import React from 'react';

function FeedbackForm({ handleChange, handleSubmit, textValue }) {
  return (
    <form onSubmit={handleSubmit} >
      <textarea
        onChange={handleChange}
        placeholder="Type feedback here."
        value={textValue}
      />
      <input type="submit" value="submit" />
    </form>
  );
}

export default FeedbackForm;
