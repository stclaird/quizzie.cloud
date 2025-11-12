import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  Text,
  View,
  TouchableOpacity,
  SafeAreaView,
  ScrollView,
  Alert,
} from 'react-native';

const API_BASE_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/';

// Debug: Log which API URL is being used
console.log('🌐 Environment:', process.env.NODE_ENV);
console.log('🔗 API Base URL:', API_BASE_URL);
console.log('🔗 EXPO_PUBLIC_API_URL:', process.env.EXPO_PUBLIC_API_URL);

export default function App() {
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [selectedAnswers, setSelectedAnswers] = useState({});
  const [score, setScore] = useState(0);
  const [showResults, setShowResults] = useState(false);
  const [loading, setLoading] = useState(true);
  const [currentView, setCurrentView] = useState('categories');
  const [answerFeedback, setAnswerFeedback] = useState({});
  const [showAnswerFeedback, setShowAnswerFeedback] = useState(false);
  const [buttonEnabled, setButtonEnabled] = useState(false);
  const [runningScore, setRunningScore] = useState(0);
  const [showQuestionResult, setShowQuestionResult] = useState(false);
  const [currentQuestionResult, setCurrentQuestionResult] = useState(null);

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      setLoading(true);
      console.log('🔍 Attempting to fetch from:', `${API_BASE_URL}categories      // Update the Quiz Screen section (around line 150-200)

      // Quiz Screen
      if (currentView === 'quiz' && questions.length > 0) {
        const currentQuestion = questions[currentQuestionIndex];
        const selectedAnswerId = selectedAnswers[currentQuestionIndex];

        // Add debug logging
        console.log('🎯 Current question data:', currentQuestion);
        console.log('🔤 Available fields:', Object.keys(currentQuestion));

        return (
          <SafeAreaView style={styles.container}>
            <View style={styles.header}>
              <TouchableOpacity onPress={selectedCategory ? backToSubcategories : backToCategories}>
                <Text style={styles.backButton}>← Back</Text>
              </TouchableOpacity>
              <Text style={styles.headerTitle}>
                {selectedCategory ? selectedCategory.category : 'Quiz'}
              </Text>
            </View>

            <View style={styles.questionHeader}>
              <Text style={styles.questionCounter}>
                Question {currentQuestionIndex + 1} of {questions.length}
              </Text>
              <Text style={styles.category}>{currentQuestion.subcategory}</Text>
            </View>

            <ScrollView style={styles.content}>
              {/* FIXED: Use questionText and add fallbacks */}
              <Text style={styles.questionText}>
                {currentQuestion.questionText || currentQuestion.text || currentQuestion.Text || 'Question text not available'}
              </Text>

              {/* Debug info - remove this later */}
              <Text style={{ fontSize: 12, color: 'gray', marginBottom: 10 }}>
                Debug: {JSON.stringify(Object.keys(currentQuestion))}
              </Text>

              <View style={styles.answersContainer}>
                {currentQuestion.answers?.map((answer) => (
                  <TouchableOpacity
                    key={answer.ID}
                    style={[
                      styles.answerButton,
                      selectedAnswerId === answer.ID && styles.selectedAnswer,
                    ]}
                    onPress={() => handleAnswerSelect(answer.ID)}
                  >
                    <Text
                      style={[
                        styles.answerText,
                        selectedAnswerId === answer.ID && styles.selectedAnswerText,
                      ]}
                    >
                      {answer.text}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </ScrollView>

            <View style={styles.footer}>
              <TouchableOpacity
                style={[
                  styles.nextButton,
                  !selectedAnswerId && styles.disabledButton,
                ]}
                onPress={nextQuestion}
                disabled={!selectedAnswerId}
              >
                <Text style={styles.nextButtonText}>
                  {currentQuestionIndex === questions.length - 1 ? 'Finish' : 'Next'}
                </Text>
              </TouchableOpacity>
            </View>
          </SafeAreaView>
        );
      }`);

      const response = await fetch(`${API_BASE_URL}categories/`);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();

      setCategories(data || []);
      console.log('✅ Categories set successfully!');
      setLoading(false);
    } catch (error) {
      console.error('❌ ERROR Details:', {
        message: error.message,
        name: error.name,
        stack: error.stack
      });
      Alert.alert('Error', `Failed to load categories: ${error.message}`);
      setLoading(false);
    }
  };

  // Add a debug button temporarily
  const DebugButton = () => (
    <TouchableOpacity
      style={[styles.primaryButton, { backgroundColor: '#ff6b6b', marginTop: 10 }]}
      onPress={() => {
        console.log('🐛 Debug button pressed');
        console.log('🌐 API URL:', API_BASE_URL);
        console.log('📊 Current categories:', categories);
        fetchCategories();
      }}
    >
      <Text style={styles.primaryButtonText}>🐛 Debug Fetch</Text>
    </TouchableOpacity>
  );

  const fetchQuestionsByCategory = async (category, subcategory = null) => {
    try {
      setLoading(true);
      let url;
      if (subcategory) {
        // Use path parameters for category and subcategory filtering
        url = `${API_BASE_URL}questions/${encodeURIComponent(category)}/${encodeURIComponent(subcategory)}`;
      } else {
        // For category-only filtering, we need to get all questions and filter client-side
        // since there's no backend endpoint for category-only filtering
        url = `${API_BASE_URL}questions/`;
      }

      console.log('🔗 Fetching questions from:', url);
      const response = await fetch(url);
      const data = await response.json();

      // Filter questions if we're doing category-only filtering
      let filteredQuestions = data || [];
      if (!subcategory && category) {
        filteredQuestions = (data || []).filter(question =>
          question.category && question.category.toLowerCase() === category.toLowerCase()
        );
        console.log(`📊 Category filtering: ${data.length} total → ${filteredQuestions.length} for category "${category}"`);
      } else if (subcategory) {
        console.log(`📊 Subcategory filtering: Got ${filteredQuestions.length} questions for "${category}/${subcategory}"`);
      }

      setQuestions(filteredQuestions);
      setCurrentQuestionIndex(0);
      setSelectedAnswers({});
      setScore(0);
      setShowResults(false);
      setCurrentView('quiz');
      setAnswerFeedback({});
      setShowAnswerFeedback(false);
      setButtonEnabled(false);
      setRunningScore(0);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching questions:', error);
      Alert.alert('Error', 'Failed to load questions');
      setLoading(false);
    }
  };

  const handleCategorySelect = (category) => {
    if (category.subcategories && category.subcategories.length > 0) {
      setSelectedCategory(category);
      setCurrentView('subcategories');
    } else {
      fetchQuestionsByCategory(category.category);
    }
  };

  const handleSubcategorySelect = (subcategory) => {
    fetchQuestionsByCategory(selectedCategory.category, subcategory.subcategoryname);
  };

  const handleAnswerSelect = (answerId) => {
    const currentSelections = selectedAnswers[currentQuestionIndex] || [];

    // Always allow multiple selections - we'll determine the question type on submission
    const selectionsArray = Array.isArray(currentSelections) ? currentSelections : [currentSelections].filter(Boolean);
    const isAlreadySelected = selectionsArray.includes(answerId);

    let newSelections;
    if (isAlreadySelected) {
      // Remove if already selected
      newSelections = selectionsArray.filter(id => id !== answerId);
    } else {
      // Add to selections
      newSelections = [...selectionsArray, answerId];
    }

    setSelectedAnswers({
      ...selectedAnswers,
      [currentQuestionIndex]: newSelections,
    });

    // Enable submit button if at least one answer is selected
    setButtonEnabled(newSelections.length > 0);
  };

  const checkAnswerAndDetermineType = async (answerId) => {
    try {
      const currentQuestion = questions[currentQuestionIndex];

      if (!currentQuestion || !currentQuestion.ID) {
        console.error('❌ Current question is invalid:', currentQuestion);
        return;
      }

      console.log('🚀 Checking answer to determine question type for question ID:', currentQuestion.ID, 'with answer ID:', answerId);
      const response = await fetch(`${API_BASE_URL}questions/answer/${currentQuestion.ID}/${answerId}`);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const answerResult = await response.json();
      const isCorrect = answerResult.iscorrect || answerResult.IsCorrect || false;
      const correctAnswers = answerResult.correctanswer || answerResult.CorrectAnswer || [];
      const totalCorrectAnswers = correctAnswers.length;
      const isMultipleChoice = totalCorrectAnswers > 1;

      console.log('🎯 Answer analysis:', {
        questionId: currentQuestion.ID,
        selectedAnswerId: answerId,
        isCorrect,
        totalCorrectAnswers,
        isMultipleChoice,
        correctAnswers
      });

      // Store the question type information and feedback
      setAnswerFeedback({
        ...answerFeedback,
        [currentQuestionIndex]: {
          isCorrect,
          selectedAnswerId: answerId,
          isMultipleChoice,
          totalCorrectAnswers,
          correctAnswers
        }
      });

      if (isMultipleChoice) {
        // This is a multiple choice question - convert to multiple selection mode
        console.log('🔄 Detected multiple choice question, switching to multi-select mode');

        // Convert the single selection to array format for multiple choice
        setSelectedAnswers({
          ...selectedAnswers,
          [currentQuestionIndex]: [answerId],
        });

        // Enable button immediately and don't show single-answer feedback
        setButtonEnabled(true);

        // Don't show the single-answer feedback for multiple choice
        setShowAnswerFeedback(false);
      } else {
        // This is a single choice question - show normal feedback
        if (isCorrect) {
          setRunningScore(prevScore => prevScore + 1);
        }

        setShowAnswerFeedback(true);

        // Enable the button after seeing feedback
        setTimeout(() => {
          setButtonEnabled(true);
        }, 1500);

        // Hide feedback after 4 seconds
        setTimeout(() => {
          setShowAnswerFeedback(false);
        }, 4000);
      }

    } catch (error) {
      console.error('Error checking answer:', error);
    }
  };

  const checkSingleAnswer = async (answerId) => {
    // This is now just a wrapper for the main checking function
    await checkAnswerAndDetermineType(answerId);
  };

  const nextQuestion = async () => {
    if (showQuestionResult) {
      // User is viewing result screen, move to next question or finish
      setShowQuestionResult(false);
      setCurrentQuestionResult(null);
      setButtonEnabled(false);

      if (currentQuestionIndex < questions.length - 1) {
        setCurrentQuestionIndex(currentQuestionIndex + 1);
      } else {
        calculateScore();
      }
    } else {
      // User is submitting their answer, check it and show result
      await submitAndShowResult();
    }
  };

  const submitAndShowResult = async () => {
    const currentQuestion = questions[currentQuestionIndex];
    const selectedAnswerData = selectedAnswers[currentQuestionIndex];

    if (!selectedAnswerData || selectedAnswerData.length === 0) {
      return;
    }

    try {
      // Convert selections to API format
      const answerString = Array.isArray(selectedAnswerData)
        ? selectedAnswerData.sort().join('')
        : selectedAnswerData.toString();

      console.log('🚀 Submitting answer:', {
        questionId: currentQuestion.ID,
        selectedAnswers: selectedAnswerData,
        answerString
      });

      const response = await fetch(`${API_BASE_URL}questions/answer/${currentQuestion.ID}/${answerString}`);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const answerResult = await response.json();
      const isCorrect = answerResult.iscorrect || answerResult.IsCorrect || false;
      const correctAnswers = answerResult.correctanswer || answerResult.CorrectAnswer || [];

      // Update running score
      if (isCorrect) {
        setRunningScore(prevScore => prevScore + 1);
      }

      // Prepare result data for display
      setCurrentQuestionResult({
        isCorrect,
        correctAnswers,
        userAnswers: selectedAnswerData,
        questionText: currentQuestion.questionText,
        allAnswers: currentQuestion.answers
      });

      // Show the result screen
      setShowQuestionResult(true);
      setButtonEnabled(true); // Enable "Continue" button

    } catch (error) {
      console.error('Error checking answer:', error);
    }
  };

  const updateRunningScoreForMultipleChoice = async () => {
    const currentQuestion = questions[currentQuestionIndex];
    const selectedAnswerData = selectedAnswers[currentQuestionIndex];

    if (!selectedAnswerData || !Array.isArray(selectedAnswerData) || selectedAnswerData.length === 0) {
      return; // No answers selected
    }

    try {
      // Convert array to string format expected by backend
      const answerString = selectedAnswerData.sort().join('');
      const response = await fetch(`${API_BASE_URL}questions/answer/${currentQuestion.ID}/${answerString}`);

      if (response.ok) {
        const answerResult = await response.json();
        const isCorrect = answerResult.iscorrect || answerResult.IsCorrect || false;

        console.log('🏃‍♂️ Running score update for multichoice:', {
          questionId: currentQuestion.ID,
          selectedAnswers: selectedAnswerData,
          isCorrect,
          currentRunningScore: runningScore
        });

        if (isCorrect) {
          setRunningScore(prevScore => prevScore + 1);
        }
      }
    } catch (error) {
      console.error('Error updating running score for multiple choice:', error);
    }
  };

  const calculateScore = async () => {
    let correctAnswers = 0;

    for (let i = 0; i < questions.length; i++) {
      const question = questions[i];
      const selectedAnswerData = selectedAnswers[i];
      const questionFeedback = answerFeedback[i];
      const isMultipleChoice = questionFeedback?.isMultipleChoice || false;

      if (isMultipleChoice && Array.isArray(selectedAnswerData)) {
        // For multiple choice, check if all selected answers are correct
        // This uses a more complex scoring logic for multiple selections
        try {
          // Convert array to string format expected by backend
          const answerString = selectedAnswerData.sort().join('');
          const response = await fetch(`${API_BASE_URL}questions/answer/${question.ID}/${answerString}`);

          if (!response.ok) {
            console.error(`Failed to check multiple choice answer for question ${question.ID}: ${response.status} ${response.statusText}`);
            continue;
          }

          const answerResult = await response.json();
          const isCorrectForScore = answerResult.iscorrect || answerResult.IsCorrect || false;

          console.log(`📊 Multiple choice score Q${i + 1}:`, {
            questionId: question.ID,
            selectedAnswers: selectedAnswerData,
            answerString,
            isCorrectForScore,
            answerResult
          });

          if (isCorrectForScore) {
            correctAnswers++;
          }
        } catch (error) {
          console.error('Error checking multiple choice answer:', error);
        }
      } else {
        // Single answer question
        const selectedAnswerId = Array.isArray(selectedAnswerData) ? selectedAnswerData[0] : selectedAnswerData;

        if (!selectedAnswerId) continue;

        try {
          const response = await fetch(`${API_BASE_URL}questions/answer/${question.ID}/${selectedAnswerId}`);

          if (!response.ok) {
            console.error(`Failed to check answer for question ${question.ID}: ${response.status} ${response.statusText}`);
            continue;
          }

          const answerResult = await response.json();
          const isCorrectForScore = answerResult.iscorrect || answerResult.IsCorrect || false;

          console.log(`📊 Single answer score Q${i + 1}:`, {
            questionId: question.ID,
            selectedAnswerId,
            isCorrectForScore,
            answerResult
          });

          if (isCorrectForScore) {
            correctAnswers++;
          }
        } catch (error) {
          console.error('Error checking single answer:', error);
        }
      }
    }

    setScore(correctAnswers);
    setShowResults(true);
  };

  const resetToCategories = () => {
    setCurrentView('categories');
    setSelectedCategory(null);
    setQuestions([]);
    setCurrentQuestionIndex(0);
    setSelectedAnswers({});
    setScore(0);
    setShowResults(false);
    setAnswerFeedback({});
    setShowAnswerFeedback(false);
    setButtonEnabled(false);
    setRunningScore(0);
  };

  const backToCategories = () => {
    setCurrentView('categories');
    setSelectedCategory(null);
  };

  const backToSubcategories = () => {
    setCurrentView('subcategories');
  };

  // Loading Screen
  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.centerContainer}>
          <Text style={styles.loadingText}>Loading...</Text>
        </View>
      </SafeAreaView>
    );
  }

  // Results Screen
  if (showResults) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={resetToCategories}>
            <Text style={styles.backButton}>← Back to Categories</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Quiz Results</Text>
        </View>
        <View style={styles.resultsContainer}>
          <Text style={styles.resultsTitle}>Quiz Complete!</Text>
          <Text style={styles.scoreText}>
            Your Score: {score} / {questions.length}
          </Text>
          <Text style={styles.percentageText}>
            {Math.round((score / questions.length) * 100)}%
          </Text>
          <TouchableOpacity style={styles.primaryButton} onPress={resetToCategories}>
            <Text style={styles.primaryButtonText}>Back to Categories</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  // Question Result Screen
  if (showQuestionResult && currentQuestionResult) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <Text style={styles.headerTitle}>Question Result</Text>
        </View>

        <View style={styles.questionHeader}>
          <Text style={styles.questionCounter}>
            Question {currentQuestionIndex + 1} of {questions.length}
          </Text>
          <View style={styles.scoreContainer}>
            <Text style={styles.runningScore}>
              Score: {runningScore} / {questions.length}
            </Text>
            <Text style={styles.percentageScore}>
              {questions.length > 0 ? Math.round((runningScore / questions.length) * 100) : 0}%
            </Text>
          </View>
        </View>

        <ScrollView style={styles.content}>
          {/* Question */}
          <Text style={styles.questionText}>{currentQuestionResult.questionText}</Text>

          {/* Overall Result */}
          <View style={[
            styles.resultContainer,
            currentQuestionResult.isCorrect ? styles.correctResult : styles.incorrectResult
          ]}>
            <Text style={styles.resultText}>
              {currentQuestionResult.isCorrect ? '✅ CORRECT!' : '❌ INCORRECT'}
            </Text>
          </View>

          {/* Answer Analysis */}
          <View style={styles.answerAnalysis}>
            <Text style={styles.analysisTitle}>Answer Analysis:</Text>

            {currentQuestionResult.allAnswers.map((answer) => {
              const isUserSelected = currentQuestionResult.userAnswers.includes(answer.ID);
              const isCorrectAnswer = currentQuestionResult.correctAnswers.some(ca => ca.ID === answer.ID);

              return (
                <View
                  key={answer.ID}
                  style={[
                    styles.analysisAnswer,
                    isUserSelected && styles.userSelectedAnswer,
                    isCorrectAnswer && styles.correctAnalysisAnswer,
                  ]}
                >
                  <Text style={[
                    styles.analysisAnswerText,
                    isUserSelected && styles.userSelectedText,
                    isCorrectAnswer && styles.correctAnalysisText,
                  ]}>
                    {answer.text}
                    {isUserSelected && ' 👆 (You selected)'}
                    {isCorrectAnswer && ' ✅ (Correct)'}
                  </Text>
                </View>
              );
            })}
          </View>
        </ScrollView>

        <View style={styles.footer}>
          <TouchableOpacity
            style={styles.nextButton}
            onPress={nextQuestion}
          >
            <Text style={styles.nextButtonText}>
              {currentQuestionIndex === questions.length - 1 ? 'Finish Quiz' : 'Continue'}
            </Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  // Quiz Screen
  if (currentView === 'quiz' && questions.length > 0) {
    const currentQuestion = questions[currentQuestionIndex];
    const selectedAnswerData = selectedAnswers[currentQuestionIndex];

    // Always treat as potentially multiple choice - user can select multiple answers
    const selectedAnswerIds = Array.isArray(selectedAnswerData) ? selectedAnswerData : (selectedAnswerData ? [selectedAnswerData] : []);

    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={selectedCategory ? backToSubcategories : backToCategories}>
            <Text style={styles.backButton}>← Back</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>
            {selectedCategory ? selectedCategory.category : 'Quiz'}
          </Text>
        </View>

        <View style={styles.questionHeader}>
          <Text style={styles.questionCounter}>
            Question {currentQuestionIndex + 1} of {questions.length}
          </Text>
          <Text style={styles.category}>{currentQuestion.subcategory}</Text>
          <View style={styles.scoreContainer}>
            <Text style={styles.runningScore}>
              Score: {runningScore} / {questions.length}
            </Text>
            <Text style={styles.percentageScore}>
              {questions.length > 0 ? Math.round((runningScore / questions.length) * 100) : 0}%
            </Text>
          </View>
        </View>

        <ScrollView style={styles.content}>
          <Text style={styles.questionText}>{currentQuestion.questionText}</Text>

          {/* No immediate feedback - show results after submission */}

          <View style={styles.answersContainer}>
            <Text style={styles.instructionText}>
              💡 Select your answer(s) and click Submit ({selectedAnswerIds.length} selected)
            </Text>
            {currentQuestion.answers?.map((answer) => {
              const feedback = answerFeedback[currentQuestionIndex];
              const isSelected = selectedAnswerIds.includes(answer.ID);

              // Check if this answer is in the list of correct answers from API
              const correctAnswerIds = feedback?.correctAnswers?.map(ca => ca.ID) || [];
              const isCorrectAnswer = correctAnswerIds.includes(answer.ID);

              // Debug logging for color coding
              if (feedback && isSelected) {
                console.log('🎨 Color coding check:', {
                  answerId: answer.ID,
                  isSelected,
                  isCorrectAnswer,
                  correctAnswerIds,
                  isMultipleChoice,
                  showCorrectFeedback,
                  showIncorrectFeedback,
                  answerText: answer.text
                });
              }

              // No immediate feedback - just show selection state
              var showCorrectFeedback = false;
              var showIncorrectFeedback = false;
              var showSelectedCorrectAnswer = false;
              var showSelectedIncorrectAnswer = false;

              return (
                <TouchableOpacity
                  key={answer.ID}
                  style={[
                    styles.answerButton,
                    isSelected && styles.selectedAnswer,
                    showCorrectFeedback && styles.correctAnswer,
                    showIncorrectFeedback && styles.incorrectAnswer,
                    showSelectedCorrectAnswer && styles.correctAnswer,
                    showSelectedIncorrectAnswer && styles.incorrectAnswer,
                  ]}
                  onPress={() => handleAnswerSelect(answer.ID)}
                >
                  <Text
                    style={[
                      styles.answerText,
                      isSelected && styles.selectedAnswerText,
                      showCorrectFeedback && styles.correctAnswerText,
                      showIncorrectFeedback && styles.incorrectAnswerText,
                      showSelectedCorrectAnswer && styles.correctAnswerText,
                      showSelectedIncorrectAnswer && styles.incorrectAnswerText,
                    ]}
                  >
                    {answer.text}
                    {(showCorrectFeedback || showSelectedCorrectAnswer) && ' ✅'}
                    {(showIncorrectFeedback || showSelectedIncorrectAnswer) && ' ❌'}
                  </Text>
                </TouchableOpacity>
              );
            })}
          </View>
        </ScrollView>

        <View style={styles.footer}>
          <TouchableOpacity
            style={[
              styles.nextButton,
              selectedAnswerIds.length === 0 && styles.disabledButton,
            ]}
            onPress={nextQuestion}
            disabled={selectedAnswerIds.length === 0}
          >
            <Text style={styles.nextButtonText}>
              Submit Answer
            </Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  // Subcategories Screen
  if (currentView === 'subcategories' && selectedCategory) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={backToCategories}>
            <Text style={styles.backButton}>← Categories</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>{selectedCategory.category}</Text>
        </View>

        <ScrollView style={styles.content}>
          <Text style={styles.sectionTitle}>Select a Subcategory:</Text>

          <View style={styles.listContainer}>
            {selectedCategory.subcategories.map((subcategory, index) => (
              <TouchableOpacity
                key={index}
                style={styles.categoryCard}
                onPress={() => handleSubcategorySelect(subcategory)}
              >
                <View style={styles.categoryContent}>
                  <Text style={styles.categoryTitle}>
                    {subcategory.subcategoryname}
                  </Text>
                  <Text style={styles.categoryArrow}>→</Text>
                </View>
              </TouchableOpacity>
            ))}
          </View>
        </ScrollView>
      </SafeAreaView>
    );
  }

  // Categories Screen (Default)
  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Quizzie Categories</Text>
      </View>

      <ScrollView style={styles.content}>
        <Text style={styles.sectionTitle}>Select a Category:</Text>

        {categories.length === 0 ? (
          <View style={styles.centerContainer}>
            <Text style={styles.errorText}>No categories available</Text>
            <TouchableOpacity style={styles.primaryButton} onPress={fetchCategories}>
              <Text style={styles.primaryButtonText}>Retry</Text>
            </TouchableOpacity>
            <DebugButton />
          </View>
        ) : (
          <View style={styles.listContainer}>
            {categories.map((category, index) => (
              <TouchableOpacity
                key={index}
                style={styles.categoryCard}
                onPress={() => handleCategorySelect(category)}
              >
                <View style={styles.categoryContent}>
                  <View>
                    <Text style={styles.categoryTitle}>{category.category}</Text>
                    <Text style={styles.subcategoryCount}>
                      {category.subcategories?.length || 0} subcategories
                    </Text>
                  </View>
                  <Text style={styles.categoryArrow}>→</Text>
                </View>
              </TouchableOpacity>
            ))}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e9ecef',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 3,
  },
  backButton: {
    fontSize: 16,
    color: '#007AFF',
    fontWeight: '600',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2c3e50',
    flex: 1,
    textAlign: 'center',
    marginLeft: -50,
  },
  content: {
    flex: 1,
    padding: 20,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 20,
  },
  listContainer: {
    gap: 15,
  },
  categoryCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 3,
  },
  categoryContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  categoryTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#2c3e50',
  },
  subcategoryCount: {
    fontSize: 14,
    color: '#6c757d',
    marginTop: 4,
  },
  categoryArrow: {
    fontSize: 20,
    color: '#007AFF',
    fontWeight: 'bold',
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    fontSize: 18,
    color: '#6c757d',
  },
  errorText: {
    fontSize: 18,
    textAlign: 'center',
    color: '#dc3545',
    marginBottom: 20,
  },
  primaryButton: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 10,
    minWidth: 120,
    alignItems: 'center',
  },
  primaryButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  questionHeader: {
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e9ecef',
  },
  questionCounter: {
    fontSize: 14,
    color: '#6c757d',
    textAlign: 'center',
  },
  category: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#495057',
    textAlign: 'center',
    marginTop: 5,
  },
  questionText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#2c3e50',
    marginBottom: 30,
    lineHeight: 26,
  },
  answersContainer: {
    gap: 15,
  },
  answerButton: {
    backgroundColor: '#fff',
    padding: 15,
    borderRadius: 10,
    borderWidth: 2,
    borderColor: '#e9ecef',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 2,
  },
  selectedAnswer: {
    borderColor: '#007AFF',
    backgroundColor: '#f0f8ff',
  },
  answerText: {
    fontSize: 16,
    color: '#495057',
  },
  selectedAnswerText: {
    color: '#007AFF',
    fontWeight: '600',
  },
  footer: {
    padding: 20,
    backgroundColor: '#fff',
    borderTopWidth: 1,
    borderTopColor: '#e9ecef',
  },
  nextButton: {
    backgroundColor: '#28a745',
    padding: 15,
    borderRadius: 10,
    alignItems: 'center',
  },
  disabledButton: {
    backgroundColor: '#6c757d',
  },
  nextButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  resultsContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  resultsTitle: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: 20,
  },
  scoreText: {
    fontSize: 24,
    color: '#495057',
    marginBottom: 10,
  },
  percentageText: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#28a745',
    marginBottom: 30,
  },
  totalQuestions: {
    fontSize: 12,
    color: '#6c757d',
    textAlign: 'center',
    marginTop: 5,
    fontStyle: 'italic',
  },
  feedbackContainer: {
    padding: 15,
    borderRadius: 10,
    marginBottom: 20,
    alignItems: 'center',
  },
  correctFeedback: {
    backgroundColor: '#d4edda',
    borderColor: '#28a745',
    borderWidth: 1,
  },
  incorrectFeedback: {
    backgroundColor: '#f8d7da',
    borderColor: '#dc3545',
    borderWidth: 1,
  },
  feedbackText: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  correctAnswer: {
    backgroundColor: '#d4edda',
    borderColor: '#28a745',
    borderWidth: 2,
    shadowColor: '#28a745',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 4,
    elevation: 4,
  },
  incorrectAnswer: {
    backgroundColor: '#f8d7da',
    borderColor: '#dc3545',
    borderWidth: 2,
    shadowColor: '#dc3545',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 4,
    elevation: 4,
  },
  correctAnswerText: {
    color: '#155724',
    fontWeight: 'bold',
    fontSize: 16,
  },
  incorrectAnswerText: {
    color: '#721c24',
    fontWeight: 'bold',
    fontSize: 16,
  },
  instructionText: {
    fontSize: 14,
    color: '#495057',
    fontWeight: '600',
    textAlign: 'center',
    marginBottom: 15,
    backgroundColor: '#fff3cd',
    padding: 12,
    borderRadius: 8,
    borderLeftWidth: 4,
    borderLeftColor: '#ffc107',
    borderWidth: 1,
    borderColor: '#ffeaa7',
  },
  scoreContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 10,
    paddingTop: 10,
    borderTopWidth: 1,
    borderTopColor: '#e9ecef',
  },
  runningScore: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#28a745',
  },
  percentageScore: {
    fontSize: 14,
    fontWeight: '600',
    color: '#007AFF',
  },
  resultContainer: {
    padding: 20,
    borderRadius: 12,
    marginBottom: 20,
    alignItems: 'center',
  },
  correctResult: {
    backgroundColor: '#d4edda',
    borderColor: '#28a745',
    borderWidth: 2,
  },
  incorrectResult: {
    backgroundColor: '#f8d7da',
    borderColor: '#dc3545',
    borderWidth: 2,
  },
  resultText: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  answerAnalysis: {
    marginTop: 20,
  },
  analysisTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 15,
    color: '#2c3e50',
  },
  analysisAnswer: {
    padding: 12,
    borderRadius: 8,
    marginBottom: 10,
    backgroundColor: '#f8f9fa',
    borderWidth: 1,
    borderColor: '#e9ecef',
  },
  userSelectedAnswer: {
    backgroundColor: '#e3f2fd',
    borderColor: '#2196f3',
  },
  correctAnalysisAnswer: {
    backgroundColor: '#e8f5e8',
    borderColor: '#4caf50',
  },
  analysisAnswerText: {
    fontSize: 16,
    color: '#495057',
  },
  userSelectedText: {
    color: '#1976d2',
    fontWeight: '600',
  },
  correctAnalysisText: {
    color: '#388e3c',
    fontWeight: '600',
  },
});
