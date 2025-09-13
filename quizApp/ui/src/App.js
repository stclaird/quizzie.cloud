import React, { useState, useEffect } from 'react';

import Category from './components/category';
import Questions from './components/questions';

const url = "http://localhost:5000/";
const categoriesUrl = [ url, 'categories/'].join('');

function App(props) {

  const [categories, setCategories] = useState([]);
  const [showCategories, setShowCategories] = useState(true);

  const [showQuestions, setShowQuestions] = useState(false);
  const [questions, setQuestions] = useState("")

  //Show Questions
  function startQuestions(catName, subCatName) {
    //Don't show categories when in question mode.
    setShowCategories(false)
    //display questions to true
    setShowQuestions(true);

    //Launch the questions
    setQuestions({
      catName : catName,
      subCatName : subCatName,
      initialQuestionIdx : 0
   })
  }

  return (
  <div className="container">
    <div className="page-header">
    </div>
    <hr></hr>
      {showCategories &&  <Category
            startQuestions={startQuestions}
            categoriesUrl={categoriesUrl}
      /> }
      {showQuestions && <Questions
            catName={questions.catName}
            subCatName={questions.subCatName}
            initialQuestionIdx={0}
      />}
  </div>
 )
}
export default App;
